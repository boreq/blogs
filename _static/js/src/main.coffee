# This function should be called on the "add photo" page.
mainAddPhoto = () ->
    initMapPickerOnForm('map-picker', 'lon', 'lat')
    initFileUploadArea('#upload-file', '.photo-form input[name=image]', '#upload-file-preview')


# This function should be called on the "add photo" page.
mainEditPhoto = () ->
    initMapPickerOnForm('map-picker', 'lon', 'lat')


# This function should be called on the "add current photo" page.
mainAddCurrentPhoto = () ->
    initLocationPreviewMap('photo-map')
    comparison = new Comparison('comparison-canvas')
    editor = new Editor('#upload-file', 'input[name=image]', 'editor-canvas',
                        historicalImageUrl)
    editor.onImageChanged = (image1, image2) ->
        comparison.setImage1(image1)
        comparison.setImage2(image2)
        comparison.redraw()
    $('#upload-button').on('click', () ->
        btn = $(this).button('loading')
        imageDataURL = editor.getEditedImage().toDataURL('image/jpeg', 0.9)
        formData = new FormData()
        formData.append('csrfmiddlewaretoken', $('input[name=csrfmiddlewaretoken]').val())
        formData.append('image', dataURItoBlob(imageDataURL))
        $.ajax(
            url: uploadAPIUrl,
            method: "POST",
            data: formData,
            processData: false,
            contentType: false
        ).done((data, textStatus, jqXHR) ->
            window.location.replace(data.url)
        ).fail((jqXHR, textStatus, errorThrown) ->
            response = $.parseJSON(jqXHR.responseText)
            if response.message?
                alert(response.message)
            btn.button('reset')
        )
    )

    initialValue = editor.getHistoricalImageAlpha() * ($(this).attr('max') - $(this).attr('min'))
    $('#editor-transparency-input').val(initialValue)
    $('#editor-transparency-input').on('input', (e) ->
        min = $(this).attr('min')
        max = $(this).attr('max')
        val = $(this).val()
        editor.setHistoricalImageAlpha(val / (max - min))
        editor.redraw()
    )


# This function should be called on the photo page.
mainPhoto = () ->
    initLocationPreviewMap('photo-map')

    loadingIndicator = $('.loading-indicator')

    loaded = 0
    imgLoaded = () ->
        loaded++
        if loaded >= 2
            loadingIndicator.hide()

    switchTo = (imgElement, autoScroll=false) ->
        currentImage = new Image()
        currentImage.onload = () ->
            comparison.setImage2(currentImage)
            imgLoaded()
            comparison.redraw()
        currentImage.src = imgElement.attr('src')
        currentPhotoId = imgElement.closest('.current-photo').attr('current_photo_id')
        $('.current-photo').removeClass('selected')
        $('.current-photo[current_photo_id=' + currentPhotoId + ']').addClass('selected')
        if autoScroll
            $('html, body').animate({
                scrollTop: $('.photo-container').offset().top
            }, 500)

    comparison = new Comparison('photo-canvas', false)

    # Load the historical image
    historicalImage = new Image()
    historicalImage.onload = () ->
        comparison.setImage1(historicalImage)
        imgLoaded()
        comparison.redraw()
    historicalImage.src = historicalImageUrl

    # Load the first image from the list of current images
    permalinkedImage = $('#permalinked-current-photo img')
    if permalinkedImage.length > 0
        switchTo(permalinkedImage.first())
    else
        images = $('#current-photos-list img')
        if images.length > 0
            switchTo(images.first())
        else
            imgLoaded()

    # Change the current photo when the user clicks a thumbnail
    $('.current-photo img').on('click', (e) ->
        loadingIndicator.show()
        switchTo($(this), true)
    )

    # Voting
    $('.photo-vote-button').on('click', (e) ->
        button = $(this)
        upvote = not button.hasClass('upvoted')
        id = button.attr('photo_id')
        $.ajax(
            url: voteAPIUrl,
            method: "POST",
            data:
                upvote: upvote
                id: id
                csrfmiddlewaretoken: $('input[name=csrfmiddlewaretoken]').val()
        ).done((data, textStatus, jqXHR) ->
            buttons = $('.photo-vote-button[photo_id=' + id + ']')
            v = parseInt(buttons.find('.score').text()[0])
            if data.upvoted
                buttons.addClass('upvoted')
                buttons.find('.score').text(v + 1)
            else
                buttons.removeClass('upvoted')
                buttons.find('.score').text(v - 1)
        ).fail((jqXHR, textStatus, errorThrown) ->
            response = $.parseJSON(jqXHR.responseText)
            if response.message?
                alert(response.message)
        )
    )

    $('.current-photo .vote-button').on('click', (e) ->
        button = $(this)
        upvote = not button.hasClass('upvoted')
        id = button.closest('.current-photo').attr('current_photo_id')
        $.ajax(
            url: currentVoteAPIUrl,
            method: "POST",
            data:
                upvote: upvote
                id: id
                csrfmiddlewaretoken: $('input[name=csrfmiddlewaretoken]').val()
        ).done((data, textStatus, jqXHR) ->
            buttons = $('.current-photo[current_photo_id=' + id + '] .vote-button')
            v = parseInt(buttons.find('.score').text()[0])
            if data.upvoted
                buttons.addClass('upvoted')
                buttons.find('.score').text(v + 1)
            else
                buttons.removeClass('upvoted')
                buttons.find('.score').text(v - 1)
        ).fail((jqXHR, textStatus, errorThrown) ->
            response = $.parseJSON(jqXHR.responseText)
            if response.message?
                alert(response.message)
        )
    )



# This function should be called on the account page.
mainAccount = () ->
    initMapPickerOnForm('map-picker', 'lon', 'lat')


# This function should be called on the profile page.
mainProfile = () ->
    target = 'profile-map'
    if $('#' + target).length
        initLocationPreviewMap(target)
    grid = $('.grid').masonry(
        itemSelector: '.grid-item'
        percentPosition: true
    )
    grid.imagesLoaded().progress(() ->
        grid.masonry('layout')
    )


# This function should be called on the explore page.
mainExplore = () ->
    initExploreMap('explore-map')
    grid = $('.grid').masonry(
        itemSelector: '.grid-item'
        percentPosition: true
    )
    grid.imagesLoaded().progress(() ->
        grid.masonry('layout')
    )

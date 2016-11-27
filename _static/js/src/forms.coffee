initFileUploadArea = (upload_area_selector, upload_input_selector, image_preview_selector) ->
    uploadArea = $(upload_area_selector)
    uploadInput = $(upload_input_selector)
    console.log(uploadInput)

    uploadArea.on('click', (e) ->
        e.preventDefault()
        uploadInput.trigger('click')
    )

    uploadInput.on('change', (e) ->
        if uploadInput[0].files and uploadInput[0].files[0]
            reader = new FileReader()

            reader.onload = (e) ->
                console.log($(image_preview_selector))
                $(image_preview_selector).attr('src', e.target.result)
                uploadArea.addClass('file-selected')

            reader.readAsDataURL(uploadInput[0].files[0])
    )

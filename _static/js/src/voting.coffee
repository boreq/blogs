urls =
    star: '/api/post/star'
    unstar: '/api/post/unstar'
    subscribe: '/api/blog/subscribe'
    unsubscribe: '/api/blog/unsubscribe'

handleVoteButtonClick = (button, urlSelect, urlUnselect, classSelected, classUnselected, successCallback, failureCallback) ->
    select = $(button).attr('isselected') == 'false'
    $(button).tooltip('destroy')
    $(button).find('i').attr('class', 'fa fa-spinner fa-spin fa-fw')

    data = {}
    key = $(button).closest('form').find('input').attr('name')
    data[key] = $(button).closest('form').find('input').val()

    $.ajax(
        url: if select then urlSelect else urlUnselect
        type: 'POST'
        data: data
    ).done((response) =>
        if select
            $(button).find('i').attr('class', classSelected)
            $(button).attr('isselected', 'true')
        else
            $(button).find('i').attr('class', classUnselected)
            $(button).attr('isselected', 'false')
        if successCallback
            successCallback(select)
    ).fail((response) =>
        $(button).find('i').attr('class', 'fa fa-times fa-fw')
        j = jQuery.parseJSON(response.responseText)
        if 'message' of j
            text = j.message
        else
            text = 'Error!'
        $(button).tooltip(
            trigger: 'manual'
            placement: 'left'
            animation: false
        ).attr('data-original-title', text)
        .tooltip('fixTitle')
        .tooltip('show')
        if failureCallback
            failureCallback(select, text)
    )

$ ->
    $('.star-form button').on('click', (e) ->
        e.preventDefault()
        handleVoteButtonClick(
            this,
            urls.star,
            urls.unstar,
            'fa fa-star fa-fw',
            'fa fa-star-o fa-fw'
        )
    )

    $('.subscribe-form-compact button').on('click', (e) ->
        e.preventDefault()
        handleVoteButtonClick(
            this,
            urls.subscribe,
            urls.unsubscribe,
            'fa fa-paper-plane fa-fw',
            'fa fa-paper-plane-o fa-fw'
        )
    )

    $('.subscribe-form button').on('click', (e) ->
        e.preventDefault()

        successCallback = (selected) =>
            if selected
                $(this).find('span').text('Unsubscribe')
                $(this).removeClass('btn-success')
                $(this).addClass('btn-danger')
            else
                $(this).find('span').text('Subscribe')
                $(this).removeClass('btn-danger')
                $(this).addClass('btn-success')

        failureCallback = (selected, errorText) =>

        handleVoteButtonClick(
            this,
            urls.subscribe,
            urls.unsubscribe,
            'fa fa-paper-plane fa-fw',
            'fa fa-paper-plane-o fa-fw',
            successCallback,
            failureCallback
        )
    )

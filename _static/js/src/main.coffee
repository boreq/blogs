createCookie = (name, value, days) ->
    if days
        date = new Date()
        date.setTime(date.getTime()+(days*24*60*60*1000))
        expires = "; expires="+date.toGMTString()
    else expires = ""
    document.cookie = name+"="+value+expires+"; path=/"

readCookie = (name) ->
    nameEQ = name + "="
    ca = document.cookie.split(';')
    for c in ca
        c = c.substring(1, c.length) while c.charAt(0) == ' '
        if c.indexOf(nameEQ) == 0
            return c.substring(nameEQ.length, c.length)
    return null

eraseCookie = (name) ->
    createCookie(name,"",-1)


$ ->
    $('time.timeago').timeago()
    $('.js-only').css('display', 'inline-block')
    $('[data-toggle="tooltip"]').tooltip(
        animation: false
    )

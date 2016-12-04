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

switchStylesheet = (alt) ->
    for stylesheet in document.styleSheets
        if stylesheet.title == "Default"
            stylesheet.disabled = alt
        if stylesheet.title == "Alternative"
            stylesheet.disabled = not alt

$ ->
    alt = readCookie("stylesheet") == "true"
    switchStylesheet(alt)
    $('#toggle-stylesheet').on('click', () ->
        alt = not alt
        if alt 
            createCookie("stylesheet", "true", 100)
        else
            createCookie("stylesheet", "false", 100)
        switchStylesheet(alt)
    )

    $('time.timeago').timeago()
    $('.js-only').css('display', 'inline-block')


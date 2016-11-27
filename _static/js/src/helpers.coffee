# Converts [x, y] on the map to [lon, lat].
coordsToLonLat = (coords) ->
    return ol.proj.transform(coords, 'EPSG:3857', 'EPSG:4326')


# Converts [lon, lat] to [x, y] on the map.
lonLatToCoords = (lonlat) ->
    return ol.proj.transform(lonlat, 'EPSG:4326', 'EPSG:3857')


# Returns the current user's location in the following format: [lon, lat].
getLocation = () ->
    if info.location?
        return [info.location.lon, info.location.lat]
    else
        return [19.938333, 50.061389]


# Calculates the distance between two points on a 2D plane.
distance = (x1, y1, x2, y2) ->
    return Math.sqrt(Math.pow(x2 - x1, 2) + Math.pow(y2 - y1, 2))


# Converts a dataURI returned by toDataURL to blob.
dataURItoBlob = (dataURI) ->
    # convert base64/URLEncoded data component to raw binary data held in a string
    if (dataURI.split(',')[0].indexOf('base64') >= 0)
        byteString = atob(dataURI.split(',')[1])
    else
        byteString = unescape(dataURI.split(',')[1])

    # separate out the mime component
    mimeString = dataURI.split(',')[0].split(':')[1].split(';')[0]

    # write the bytes of the string to a typed array
    ia = new Uint8Array(byteString.length)
    for i in [0...byteString.length]
        ia[i] = byteString.charCodeAt(i)

    return new Blob([ia], {type:mimeString})

mercatorOrigin = [-20037508.342789244, 20037508.342789244]
mercatorWorldSize = mercatorOrigin[1] * 2


markerStyle = () ->
    new ol.style.Style
        image: new ol.style.Icon
            src: '/static/gallery/img/marker.png'
            anchor: [0.5, 1]


altMarkerStyle = () ->
    new ol.style.Style
        image: new ol.style.Icon
            src: '/static/gallery/img/marker_alt.png'
            anchor: [0.5, 1]


# Create a tiled ol.layer.Vector with ol.source.TileVector.
# extract: (reponse) -> [{lat, lon, ...}, ...]
# render: (coords, {lat, lon, ...}) -> ol.Feature
makeTiledLayer = (url, style, maxResolution, extract, render) ->
    downloaded = []
    isDownloaded = (url) -> url in downloaded
    setDownloaded = (url) -> downloaded.push(url)

    source = new ol.source.TileVector
        url: url
        tileGrid: new ol.tilegrid.TileGrid
            minZoom: 0
            extent: ol.proj.get('EPSG:3857').getExtent()
            origin: mercatorOrigin
            resolutions: (mercatorWorldSize / (256 * Math.pow(2, zoom)) for zoom in [0..21])
        projection: 'EPSG:3857',
        tileLoadFunction: (url, callback) ->
            if isDownloaded(url)
                return
            $.ajax(
                url: url,
                type: 'GET'
            ).done (response) ->
                if isDownloaded(url)
                    return
                setDownloaded(url)
                features = []
                data = extract(response)
                for v in data
                    c = [v.lon, v.lat]
                    cord = ol.proj.fromLonLat(c)
                    if render == null || render == undefined
                        feature = new ol.Feature
                            geometry: new ol.geom.Point(cord)
                            data: v
                    else
                        feature = render(cord, v)
                    if feature != null
                        features.push(feature)
                callback(features)

    new ol.layer.Vector
        source: source
        style: style
        maxResolution: maxResolution


initMap = (target) ->
    worldLayer = new ol.layer.Tile
        source: new ol.source.OSM

    photosLayer = makeTiledLayer(
        urls.photosTilesAPI
        undefined
        2000
        (response) ->
            response.photos
        (cord, v) ->
            feature = new ol.Feature
                geometry: new ol.geom.Point(cord)
                data: v
            feature.setStyle(if v.current_photos > 0 then markerStyle() else altMarkerStyle())
            return feature
    )

    # Restore the previous position of the map if available
    hash = window.location.hash.substring(1)
    params = hash.split(',')
    if params.length == 3
        location = [parseFloat(params[0]), parseFloat(params[1])]
        zoom = parseInt(params[2])
    else
        location = getLocation()
        zoom = 14

    map = new ol.Map
        target: target
        layers: [worldLayer, photosLayer]
        controls: new ol.Collection()
        view: new ol.View
            center: lonLatToCoords(location)
            zoom: zoom

    # Redirect to a photo view after a photo is clicked
    selectInteraction = new ol.interaction.Select
        condition: ol.events.condition.singleClick
        layers: [photosLayer]
        multi: false
    selectInteraction.on('select', (e) ->
        for f in e.selected
            data = f.get('data')
            window.location = data.url
    )
    map.addInteraction(selectInteraction)

    # Store the position and zoom of the map in the url
    map.on('moveend', (e) ->
        lonLat = coordsToLonLat(map.getView().getCenter())
        zoom = map.getView().getZoom()
        anchor = '#' + lonLat[0] + ',' + lonLat[1] + ',' + zoom
        history.replaceState(undefined, undefined, anchor)
    )


initMapPickerOnForm = (target, lonInputName, latInputName) ->
    callbackOnSelection = (lonlat) ->
        $("input[name=#{lonInputName}]").val(lonlat[0])
        $("input[name=#{latInputName}]").val(lonlat[1])

    lon = $("input[name=#{lonInputName}]").val()
    lat = $("input[name=#{latInputName}]").val()
    if lon and lat
        initial = [parseFloat(lon), parseFloat(lat)]
    else
        initial = null

    initMapPicker(target, callbackOnSelection, initial)


# Creates a map which is used to select a location. When the map is clicked
# a marker is placed on it and the coordinates of the selected position are
# passed to the provided callback function.
#
# target: the id of a div to be used as a map container.
# callback: a function which will be passed a [lon, lat] array as an argument.
# initial: the initial position of the marker as [lon, lat] or null.
initMapPicker = (target, callback, initial) ->
    worldSource = new ol.source.OSM
    worldLayer = new ol.layer.Tile
        source: worldSource

    markersSource = new ol.source.Vector({})
    markersLayer = new ol.layer.Vector
        source: markersSource

    map = new ol.Map
        target: target
        layers: [worldLayer, markersLayer]
        controls: new ol.Collection()
        view: new ol.View
            center: lonLatToCoords(getLocation())
            zoom: 14

    interaction = new ol.interaction.Select
        condition: ol.events.condition.singleClick
        layers: [worldLayer]
        multi: false

    map.addInteraction(interaction)

    placeMarker = (coords) ->
        feature = new ol.Feature({})
        feature.setGeometry(new ol.geom.Point(coords))
        feature.setStyle(markerStyle())
        markersSource.clear()
        markersSource.addFeature(feature)

    if initial
        coords = lonLatToCoords(initial)
        placeMarker(coords)
        map.getView().setCenter(coords)

    interaction.on('select', (e) ->
        lonlat = coordsToLonLat(e.mapBrowserEvent.coordinate)
        callback(lonlat)
        placeMarker(e.mapBrowserEvent.coordinate)
    )


initLocationPreviewMap = (target) ->
    lon = parseFloat($('#' + target).attr('lon'))
    lat = parseFloat($('#' + target).attr('lat'))
    coords = lonLatToCoords([lon, lat])

    worldSource = new ol.source.OSM
    worldLayer = new ol.layer.Tile
        source: worldSource

    markersSource = new ol.source.Vector({})
    markersLayer = new ol.layer.Vector
        source: markersSource

    map = new ol.Map
        target: target
        layers: [worldLayer, markersLayer]
        controls: new ol.Collection()
        view: new ol.View
            center: coords
            zoom: 14

    # Add a marker
    feature = new ol.Feature({})
    feature.setGeometry(new ol.geom.Point(coords))
    feature.setStyle(markerStyle())
    markersSource.addFeature(feature)


initExploreMap = (target) ->
    lon = parseFloat($('#' + target).attr('lon'))
    lat = parseFloat($('#' + target).attr('lat'))
    coords = lonLatToCoords([lon, lat])

    worldSource = new ol.source.OSM
    worldLayer = new ol.layer.Tile
        source: worldSource

    markersSource = new ol.source.Vector({})
    markersLayer = new ol.layer.Vector
        source: markersSource

    map = new ol.Map
        target: target
        layers: [worldLayer, markersLayer]
        controls: new ol.Collection()
        view: new ol.View
            center: coords
            zoom: 14

    interaction = new ol.interaction.Select
        condition: ol.events.condition.singleClick
        layers: [worldLayer]
        multi: false

    map.addInteraction(interaction)

    placeMarker = (coords) ->
        feature = new ol.Feature({})
        feature.setGeometry(new ol.geom.Point(coords))
        feature.setStyle(markerStyle())
        markersSource.clear()
        markersSource.addFeature(feature)

    placeMarker(coords)

    interaction.on('select', (e) ->
        lonlat = coordsToLonLat(e.mapBrowserEvent.coordinate)
        placeMarker(e.mapBrowserEvent.coordinate)
        url = urls.explore.replace('/0/', '/' + lonlat[0] + '/')
        url = url.replace('/1/', '/' + lonlat[1] + '/')
        window.location = url
    )

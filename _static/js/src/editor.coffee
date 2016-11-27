class Bounds

    constructor: (@canvas) ->
        @x1 = null
        @y1 = null
        @x2 = null
        @y2 = null

    getX1: () =>
        return @x1 * @canvas.width

    getY1: () =>
        return @y1 * @canvas.height

    getX2: () =>
        return @x2 * @canvas.width

    getY2: () =>
        return @y2 * @canvas.height

    setX1: (x1) =>
        @x1 = x1 / @canvas.width

    setY1: (y1) =>
        @y1 = y1 / @canvas.height

    setX2: (x2) =>
        @x2 = x2 / @canvas.width

    setY2: (y2) =>
        @y2 = y2 / @canvas.height


class Editor

    @CORNER_TOP_LEFT: 0
    @CORNER_TOP_RIGHT: 1
    @CORNER_BOTTOM_LEFT: 2
    @CORNER_BOTTOM_RIGHT: 3
    @IMAGE: 4

    constructor: (upload_area_selector, upload_input_selector, canvas_id, historicalImageUrl) ->
        @canvas = document.getElementById(canvas_id)
        @ctx = @canvas.getContext('2d')

        @historicalImage = new Image()
        @historicalImage.src = historicalImageUrl
        @currentImage = null

        @dragging = false
        @draggedPoint = null
        @dragOffset = null
        @bounds = null

        @pointSelectionDistance = 10 # px
        @historicalImageAlpha = 0.5

        @uploadArea = $(upload_area_selector)
        @uploadInput = $(upload_input_selector)
        @uploadArea.on('click', @onUploadAreaClick)
        @uploadInput.on('change', @onInputChange)

        $('body').on('mousedown', @onMouseDown)
        $('body').on('mouseup', @onMouseUp)
        $('body').on('mousemove', @onMouseMove)
        $(window).on('resize', () => @redraw())

        @tmpCanvas = document.createElement('canvas')
        @tmpCanvasCtx = @tmpCanvas.getContext('2d')

    onUploadAreaClick: (e) =>
        e.preventDefault()
        @uploadInput.trigger('click')

    onInputChange: (e) =>
        if @uploadInput[0].files and @uploadInput[0].files[0]
            img = new Image()
            img.onload = () =>
                @currentImage = img
                @bounds = null
                @uploadArea.addClass('collapsed')
                $('.current-photo-step').removeClass('hidden')
                $('.current-photo-tip').slideUp(1000)
                @redraw()
            img.src = URL.createObjectURL(e.target.files[0])

    onMouseDown: (e) =>
        if @bounds != null
            @dragging = true
            @selectDraggedPoint(@mouseX, @mouseY)

    onMouseUp: (e) =>
        @dragging = false
        @draggedPoint = null
        @dragOffset = null

    onMouseMove: (e) =>
        [@mouseX, @mouseY] = @calculateMousePosition(e)
        if @dragging and @draggedPoint != null
            @dragDraggedPoint(@mouseX, @mouseY)

    calculateMousePosition: (e) =>
        rect = @canvas.getBoundingClientRect()
        return [e.clientX - rect.left, e.clientY - rect.top]

    selectDraggedPoint: (mouseX, mouseY) =>
        @draggedPoint = null

        # The fact that the code for this point is defined first DOES matter
        if mouseX > @bounds.getX1() and mouseX < @bounds.getX2() and mouseY > @bounds.getY1() and mouseY < @bounds.getY2()
            if @draggedPoint != 4 or @dragOffset == null
                @dragOffset =
                    x: ((@bounds.getX1() + @bounds.getX2()) / 2) - mouseX
                    y: ((@bounds.getY1() + @bounds.getY2()) / 2) - mouseY
            @draggedPoint = Editor.IMAGE

        if distance(mouseX, mouseY, @bounds.getX1(), @bounds.getY1()) < @pointSelectionDistance
            @draggedPoint = Editor.CORNER_TOP_LEFT

        if distance(mouseX, mouseY, @bounds.getX2(), @bounds.getY1()) < @pointSelectionDistance
            @draggedPoint = Editor.CORNER_TOP_RIGHT

        if distance(mouseX, mouseY, @bounds.getX1(), @bounds.getY2()) < @pointSelectionDistance
            @draggedPoint = Editor.CORNER_BOTTOM_LEFT

        if distance(mouseX, mouseY, @bounds.getX2(), @bounds.getY2()) < @pointSelectionDistance
            @draggedPoint = Editor.CORNER_BOTTOM_RIGHT

    dragDraggedPoint: (mouseX, mouseY) =>
        historicalRatio = @historicalImage.width / @historicalImage.height

        # Resize
        if @draggedPoint == Editor.CORNER_TOP_LEFT
            newX1 = mouseX
            if newX1 < 0
                newX1 = 0
            newY1 = @bounds.getY2() - ((@bounds.getX2() - newX1) * 1/historicalRatio)
            if newY1 < 0
                newY1 = 0
                newX1 = @bounds.getX2() - (historicalRatio * (@bounds.getY2() - newY1))
            @bounds.setX1(newX1)
            @bounds.setY1(newY1)
            @redraw()

        if @draggedPoint == Editor.CORNER_TOP_RIGHT
            newX2 = mouseX
            if newX2 >= @canvas.width
                newX2 = @canvas.width - 1
            newY1 = @bounds.getY2() - ((newX2 - @bounds.getX1()) * 1/historicalRatio)
            if newY1 < 0
                newY1 = 0
                newX2 = @bounds.getX1() + (historicalRatio * (@bounds.getY2() - newY1))
            @bounds.setX2(newX2)
            @bounds.setY1(newY1)
            @redraw()

        if @draggedPoint == Editor.CORNER_BOTTOM_LEFT
            newX1 = mouseX
            if newX1 < 0
                newX1 = 0
            newY2 = @bounds.getY1() + ((@bounds.getX2() - newX1) * 1/historicalRatio)
            if newY2 >= @canvas.height
                newY2 = @canvas.height - 1
                newX1 = @bounds.getX2() - (historicalRatio * (newY2 - @bounds.getY1()))
            @bounds.setX1(newX1)
            @bounds.setY2(newY2)
            @redraw()

        if @draggedPoint == Editor.CORNER_BOTTOM_RIGHT
            newX2 = mouseX
            if newX2 >= @canvas.width
                newX2 = @canvas.width - 1
            newY2 = @bounds.getY1() + ((newX2 - @bounds.getX1()) * 1/historicalRatio)
            if newY2 >= @canvas.height
                newY2 = @canvas.height - 1
                newX2 = @bounds.getX1() + (historicalRatio * (newY2 - @bounds.getY1()))
            @bounds.setX2(newX2)
            @bounds.setY2(newY2)
            @redraw()


        # Move
        if @draggedPoint == Editor.IMAGE
            newX1 = mouseX + @dragOffset.x - ((@bounds.getX2() - @bounds.getX1()) / 2)
            newY1 = mouseY + @dragOffset.y - ((@bounds.getY2() - @bounds.getY1()) / 2)
            newX2 = mouseX + @dragOffset.x + ((@bounds.getX2() - @bounds.getX1()) / 2)
            newY2 = mouseY + @dragOffset.y + ((@bounds.getY2() - @bounds.getY1()) / 2)
            if newX1 < 0
                newX1 = 0
                newX2 = newX1 + (@bounds.getX2() - @bounds.getX1())
            if newX2 >= @canvas.width
                newX2 = @canvas.width - 1
                newX1 = newX2 - (@bounds.getX2() - @bounds.getX1())
            if newY1 < 0
                newY1 = 0
                newY2 = newY1 + (@bounds.getY2() - @bounds.getY1())
            if newY2 >= @canvas.height
                newY2 = @canvas.height - 1
                newY1 = newY2 - (@bounds.getY2() - @bounds.getY1())
            @bounds.setX1(newX1)
            @bounds.setY1(newY1)
            @bounds.setX2(newX2)
            @bounds.setY2(newY2)
            @redraw()

    redraw:  () =>
        if not @currentImage
            return

        # Scale the canvas
        width = Math.min($(@canvas).parent().width(), @currentImage.width)
        height = width * (@currentImage.height / @currentImage.width)
        $(@canvas).attr('width', width)
        $(@canvas).attr('height', height)

        # Draw the current image
        @ctx.globalAlpha = 1
        @ctx.drawImage(@currentImage, 0, 0, width, height)

        # If this is the first time this image is being drawn pick the initial bounds
        if @bounds == null
            @bounds = new Bounds(@canvas)
            currentRatio = width / height
            historicalRatio = @historicalImage.width / @historicalImage.height
            if currentRatio < historicalRatio
                imgWidth = 0.8 * width
                imgHeight = 0.8 * width * 1/historicalRatio
                @bounds.setX1(0.1 * width)
                @bounds.setX2(0.1 * width + imgWidth)
                @bounds.setY1(0.1 * width)
                @bounds.setY2(0.1 * width + imgHeight)
            else
                imgWidth = 0.8 * height * historicalRatio
                imgHeight = 0.8 * height
                @bounds.setX1(0.1 * height)
                @bounds.setX2(0.1 * height + imgWidth)
                @bounds.setY1(0.1 * height)
                @bounds.setY2(0.1 * height + imgHeight)

        # Draw the historical image
        @ctx.globalAlpha = @historicalImageAlpha
        @ctx.drawImage(@historicalImage, @bounds.getX1(), @bounds.getY1(),
                       @bounds.getX2() - @bounds.getX1(), @bounds.getY2() - @bounds.getY1())

        # Drag markers and photo borders
        @ctx.globalAlpha = 1

        @drawPhotoBorder(@bounds.getX1(), @bounds.getY1(), @bounds.getX2(), @bounds.getY1())
        @drawPhotoBorder(@bounds.getX2(), @bounds.getY1(), @bounds.getX2(), @bounds.getY2())
        @drawPhotoBorder(@bounds.getX2(), @bounds.getY2(), @bounds.getX1(), @bounds.getY2())
        @drawPhotoBorder(@bounds.getX1(), @bounds.getY2(), @bounds.getX1(), @bounds.getY1())

        @drawDragMarker(@bounds.getX1(), @bounds.getY1())
        @drawDragMarker(@bounds.getX2(), @bounds.getY1())
        @drawDragMarker(@bounds.getX1(), @bounds.getY2())
        @drawDragMarker(@bounds.getX2(), @bounds.getY2())

        @handleOnImageChanged()

    handleOnImageChanged: () =>
        if not @bounds || not @onImageChanged
            return
        @onImageChanged(@historicalImage, @getEditedImage())

    getEditedImage: () =>
        if not @bounds
            return null
        normalizedWidth = @bounds.x2 - @bounds.x1
        normalizedHeight = @bounds.y2 - @bounds.y1
        @tmpCanvas.width = normalizedWidth * @currentImage.width
        @tmpCanvas.height = normalizedHeight * @currentImage.height
        @tmpCanvasCtx.drawImage(
            @currentImage,
            @bounds.x1 * @currentImage.width,
            @bounds.y1 * @currentImage.height,
            @tmpCanvas.width,
            @tmpCanvas.height,
            0,
            0,
            @tmpCanvas.width,
            @tmpCanvas.height
        )
        return @tmpCanvas

    drawDragMarker:  (x, y) =>
        @ctx.fillStyle = 'rgb(228, 45, 45)'
        @ctx.beginPath()
        @ctx.arc(x, y, @pointSelectionDistance, 0, 2 * Math.PI, false)
        @ctx.fill()

    drawPhotoBorder: (x1, y1, x2, y2) =>
        @ctx.strokeStyle = 'rgb(228, 45, 45)'
        @ctx.beginPath()
        @ctx.moveTo(x1, y1)
        @ctx.lineTo(x2, y2)
        @ctx.stroke()

    setHistoricalImageAlpha: (alpha) =>
        @historicalImageAlpha = alpha

    getHistoricalImageAlpha: () =>
        return @historicalImageAlpha

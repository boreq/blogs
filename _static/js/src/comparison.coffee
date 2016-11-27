class Comparison

    @LINE_COLOR: 'rgb(228, 45, 45)'
    @X_OFFSET: 50

    constructor: (canvas_id, @requireBothImages=true) ->
        @canvas = document.getElementById(canvas_id)
        @ctx = @canvas.getContext('2d')
        @linePosition = 0.5

        $('body').on('mousemove', @onMouseMove)
        $(window).on('resize', @onWindowResize)

    onMouseMove: (e) =>
        if not @image1 or not @image2
            return

        [mouseX, mouseY] = @calculateMousePosition(e)
        if mouseX > -Comparison.X_OFFSET and mouseX < @canvas.width + Comparison.X_OFFSET and
           mouseY > 0 and mouseY < @canvas.height
            @linePosition = Math.min(Math.max(mouseX / @canvas.width, 0), 1)
            @redraw()

    onWindowResize: (e) =>
        @redraw()

    calculateMousePosition: (e) =>
        rect = @canvas.getBoundingClientRect()
        return [e.clientX - rect.left, e.clientY - rect.top]

    redraw:  () =>
        if not @image1
            return

        if @requireBothImages and not @image2
            return

        image1Ratio = @image1.width / @image1.height
        if @image2
            image2Ratio = @image2.width / @image2.height
            if Math.abs(1 - image1Ratio / image2Ratio) > 0.005
                console.log('Comparison: Warning, the image ratios are different! image1Ratio/image2Ratio=', image1Ratio / image2Ratio)

        # Scale the canvas
        canvasWidth = Math.min($(@canvas).parent().width(), @image1.width)
        canvasHeight = canvasWidth * 1/image1Ratio
        $(@canvas).attr('width', canvasWidth)
        $(@canvas).attr('height', canvasHeight)

        # Draw image1
        @ctx.drawImage(@image1, 0, 0, @canvas.width, @canvas.height)

        # Draw image2
        if @image2
            @ctx.drawImage(@image2,
                           @linePosition * @image2.width, 0, (1 - @linePosition) * @image2.width, @image2.height,
                           @linePosition * @canvas.width, 0, (1 - @linePosition) * @canvas.width, @canvas.height)

            # Draw a line between images
            @ctx.strokeStyle = Comparison.LINE_COLOR
            @ctx.beginPath()
            @ctx.moveTo(@linePosition * @canvas.width, 0)
            @ctx.lineTo(@linePosition * @canvas.width, @canvas.height)
            @ctx.stroke()

    setImage1: (image) =>
        @image1 = image

    setImage2: (image) =>
        @image2 = image

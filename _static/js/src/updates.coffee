createDataSet = (backgroundColor, borderColor, label, data) ->
    {
        label: label
        fill: false
        lineTension: 0
        backgroundColor: backgroundColor
        borderColor: borderColor
        pointBorderColor: borderColor
        pointBackgroundColor: borderColor
        pointBorderWidth: 1
        pointHoverRadius: 5
        pointHoverBackgroundColor: "rgba(75,192,192,1)"
        pointHoverBorderColor: "rgba(220,220,220,1)"
        pointHoverBorderWidth: 2
        pointRadius: 5
        pointHitRadius: 10
        data: data
    }

createChart = (data) ->
    ctx = document.getElementById('updates-chart')
    data = {
        labels: data.labels,
        datasets: [
            createDataSet('rgba(92, 184, 92, 0.4)', 'rgba(92, 184, 92, 1)', 'Succeeded', data.success),
            createDataSet('rgba(217, 83, 79, 0.4)', 'rgba(217, 83, 79, 1)', 'Failed', data.failure),
        ]
    }
    myLineChart = new Chart(ctx,
        type: 'line'
        data: data
        options:
            legend:
                display: false
            scales:
                xAxes: [
                    type: 'time'
                    time:
                        unit: 'day'
                        displayFormats:
                            day: 'YYYY-MM-DD'
                ]
    )

$ ->
    if not $('#updates-chart').length
        return

    $.ajax(
        url: '/api/updates/chart.json'
        type: 'GET'
    ).done((response) =>
        createChart(response)
    ).fail((response) =>
        console.log('Error', response)
    )


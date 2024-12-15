let myChart;  
let testChart;
let prec = 2

function switchTab(tab) {
        document.querySelectorAll('.tab-content').forEach(tabContent => {
            tabContent.classList.remove('active');
        });
    
        document.querySelectorAll('.tabs div').forEach(tabElement => {
            tabElement.classList.remove('active');
        });
    
        document.getElementById('tab-' + tab).classList.add('active');
        if (tab === 'plot') {
            document.querySelector('.plot-content').classList.add('active');
        } else if (tab === 'table') {
            document.querySelector('.table-content').classList.add('active');
        } else if (tab === 'best-result') {
            document.querySelector('.best-result-content').classList.add('active');
        }
    }


async function algTest() {
    try {
        const response = await fetch("/alg_test", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({
                "data": "test"
            })
        });

        if (response.ok) {
            const data = await response.json(); 
            console.log(data); 

            generateTestChart(data);
            generateStatsTable(data);
        } else {
            console.error("Błąd przy pobieraniu danych");
        }
    } catch (error) {
        console.error("Błąd przy pobieraniu danych", error);
    }
}

function generateStatsTable(data) {
    const statsMap = data.stats_map;  
    const percentages = data.percentages;
    const tableBody = document.getElementById('test-table-body');

    tableBody.innerHTML = '';

    const row = document.createElement("tr");

    const fullWidthCell = document.createElement("td");
    fullWidthCell.colSpan = 5;
    const successPercent = (data.success_sum / 1000) * 100;
    fullWidthCell.innerText = `Spośród 1000 prób sukces osiągnęło ${data.success_sum}, co daje ${successPercent.toFixed(2)}% skuteczności programu.`;

    row.appendChild(fullWidthCell);
    tableBody.appendChild(row);

    // Convert statsMap to an array of entries and sort based on success count
    const sortedStats = Object.entries(statsMap).sort((a, b) => b[1] - a[1]);

    for (let [iteration, successCount] of sortedStats) {
        const row = document.createElement("tr");
        const iterationCell = document.createElement("td");
        const successCountCell = document.createElement("td");
        const percentageCell = document.createElement("td");

        iterationCell.textContent = iteration;
        successCountCell.textContent = successCount;
        percentageCell.textContent = (percentages[iteration] * 100).toFixed(3);

        row.appendChild(iterationCell);
        row.appendChild(percentageCell);
        row.appendChild(successCountCell);
        tableBody.appendChild(row);
    }
}


function generateTestChart(data) {
    const percentages = data.percentages;  
    const iterations = Object.keys(percentages); 
    const values = Object.values(percentages).map(p => p * 100); 

    if (testChart) {
        testChart.destroy();
    }


    const ctx = document.getElementById("test-plot").getContext("2d");
    testChart = new Chart(ctx, {
        type: 'line',  
        data: {
            labels: iterations,  
            datasets: [{
                label: 'Procent sukcesów',
                data: values,  
                borderColor: 'rgba(75, 192, 192, 1)', 
                fill: false,  
                tension: 0.1  
            }]
        },
        options: {
            responsive: true,
            scales: {
                x: {
                    title: {
                        display: true,
                        text: 'Iteracja'
                    }
                },
                y: {
                    title: {
                        display: true,
                        text: 'Procent sukcesów (%)'
                    },
                    ticks: {
                        beginAtZero: true,
                        max: 100, 
                        stepSize: 5  
                    }
                }
            }
        }
    });
}



async function calculate() {
    const a = parseFloat(document.getElementById("a").value);
    const b = parseFloat(document.getElementById("b").value);
    const d = parseFloat(document.getElementById("precision").value);
    const T = parseInt(document.getElementById("T").value);

    prec = d.toString().length - 2;

    const requestData = { a: a, b: b, d: d, T: T };

    try {
        const response = await fetch("/calculate", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(requestData),
        });
    
        if (response.ok) {
            const data = await response.json();
            console.log(data);
    
            if (myChart) {
                myChart.destroy();
            }
    
            generateChart(data);
    
            const tableBody = document.getElementById("table-body");
            tableBody.innerHTML = ""; 
    
            data.vc_data.forEach((vcData, index) => {
                const vcs = vcData.Vcs.map(vc => vc.x_real.toFixed(prec)).join(", ");
                const bins = vcData.Vcs.map(vc => vc.x_bin).join(", ");
                const fxs = vcData.Vcs.map(vc => vc.fx.toFixed(6)).join(", ");
                const maxFx = data.max_results[index] ? data.max_results[index].MaxFx : "N/A";
    
                const row = document.createElement("tr");
                row.innerHTML = `
                    <td>${index + 1}</td>
                    <td>${vcs}</td>
                    <td>${bins}</td>
                    <td>${fxs}</td>
                    <td>${maxFx}</td>
                `;
    
                tableBody.appendChild(row);
            });
    
        
        }
    } catch (error) {
        console.error("Błąd przy pobieraniu danych", error);
    }
    
}

function generateChart(data) {
    const maxResults = data.max_results;
    const iterData = data.vc_data;

    const maxFxData = maxResults.map(maxStep => ({
        x: maxStep.T,
        y: maxStep.MaxFx,
    }));

    const stepsByIteration = [];
    iterData.forEach((iter, genIndex) => {
        const stepsLength = iter.Steps.length
        const stepDivide = 1.0 / (stepsLength - 1)
        const steps = iter.Steps.map((step, stepIndex) => ({
            x: genIndex + (stepIndex) * stepDivide,
            y: step.Fx,
        }));
        stepsByIteration.push(steps);
    });

    const ctx = document.getElementById("myChart").getContext("2d");

    const IterationStepsDatasets = stepsByIteration.map((steps, index) => ({
        label: `Vc przy t = ${index+1}`,
        data: steps,
        borderColor: "rgba(0, 0, 255, 0.5)",
        showLine: true,
        fill: false,
        tension: 0.1,
        pointRadius: 2,
    }));

    const allDatasets = [
        {
            label: "MaxFx dla iteracji",
            data: maxFxData,
            borderColor: "red",
            fill: false,
            showLine: true,
            tension: 0.1,
        },
        ...IterationStepsDatasets,
    ];

    myChart = new Chart(ctx, {
        type: "scatter",
        data: {
            datasets: allDatasets,
        },
        options: {
            responsive: true,
            plugins: {
                legend: {
                    position: "top",
                    labels: {
                        filter: function (legendItem, chartData) {
                            return legendItem.text === "MaxFx dla iteracji";
                        },
                    },
                },
                tooltip: {
                    callbacks: {
                        label: function (context) {
                            const { dataset, raw } = context;
                            return `${dataset.label}: (${raw.x}, ${raw.y})`;
                        },
                    },
                },
            },
            scales: {
                x: {
                    title: {
                        display: true,
                        text: "Iteracje (T)",
                    },
                    ticks: {
                        stepSize: 1, 
                    },
                },
                y: {
                    title: {
                        display: true,
                        text: "Wartości funkcji",
                    },
                },
            },
        },
    });
}


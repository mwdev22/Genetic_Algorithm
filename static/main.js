let myChart;  
let prec = 2

// function startListeningForResults() {
//     const resultsTableBody = document.getElementById("results-table").querySelector("tbody");

//     const eventSource = new EventSource("/alg_test");

//     eventSource.onmessage = function(event) {
//         const data = JSON.parse(event.data);

//         console.log("Received data:", data);

//         const row = document.createElement("tr");
//         row.innerHTML = `
//             <td>${data.N}</td>
//             <td>${data.T}</td>
//             <td>${data.f_avg}</td>
//         `;

//         resultsTableBody.appendChild(row);

//         sortTableByFAvg();
//     };

//     function sortTableByFAvg() {
//         const rows = Array.from(resultsTableBody.querySelectorAll("tr"));
        
//         rows.sort((a, b) => {
//             const favgA = parseFloat(a.cells[4].textContent);
//             const favgB = parseFloat(b.cells[4].textContent);
//             return favgB - favgA;
//         });

//         rows.forEach(row => resultsTableBody.appendChild(row));
//     }

//     eventSource.onerror = function(err) {
//         console.error("EventSource failed:", err);
//         eventSource.close();
//     };
// };


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


// async function algTest() {

//     try {
//         const response = await fetch("/alg_test", {
//             method: "POST",
//             headers: {
//                 "Content-Type": "application/json"
//             },
//             body: JSON.stringify({
//                 "data":"test"
//             })
//         });
//         if (response.ok) {
//             const data = await response.text();
//             console.log(data)
//             startListeningForResults();
//         } else {
//             console.error("Błąd przy pobieraniu danych");
//         }
//     } catch (error) {
//         console.error("Błąd przy pobieraniu danych", error);
//     }
    
// }

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
        label: `Vc przy t = ${index}`,
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
                        text: "Iteration (T)",
                    },
                    ticks: {
                        stepSize: 1, 
                    },
                },
                y: {
                    title: {
                        display: true,
                        text: "Fx Values",
                    },
                },
            },
        },
    });
}

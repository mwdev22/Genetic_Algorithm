let myChart;  
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
    } else {
        document.querySelector('.table-content').classList.add('active');
    }
}

async function calculate() {
    const a = parseFloat(document.getElementById("a").value);
    const b = parseFloat(document.getElementById("b").value);
    const d = parseFloat(document.getElementById("precision").value);
    const N = parseInt(document.getElementById("N").value);
    const T = parseInt(document.getElementById("T").value);
    const pk = parseFloat(document.getElementById("pk").value);
    const pm = parseFloat(document.getElementById("pm").value);
    const elite = document.getElementById("elite").checked;

    prec = d.toString().length - 2

    const requestData = { a: a, b: b, d: d, T: T, N: N, pk: pk, pm: pm, elite: elite };



    const tableBody = document.getElementById("table-body");
    tableBody.innerHTML = "<tr><td colspan='5'>Ładowanie...</td></tr>";

    try {
        const response = await fetch("/calculate", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify(requestData)
        });

        if (response.ok) {
            const data = await response.json();
            console.log(data)
            const finalGenData = data.results;
            tableBody.innerHTML = '';  
            finalGenData.sort((a, b) => b.percent - a.percent);

            finalGenData.forEach(result => {
                const row = document.createElement('tr');
                row.innerHTML = `
                    <td>${result.x_real.toFixed(prec)}</td>
                    <td>${result.x_bin}</td>
                    <td>${result.fx}</td>
                    <td>${result.percent.toFixed(prec)}</td>
                    <td>${result.count}</td>
                `;
                tableBody.appendChild(row);
            });

            const genStats = data.gen_stats;
            const fmin = [];
            const fmax = [];
            const favg = [];

            genStats.forEach(stat => {
                fmin.push(stat.f_min);
                fmax.push(stat.f_max);
                favg.push(stat.f_avg);
            });

            if (myChart) {
                myChart.destroy();
            }

            const chartData = {
                labels: Array.from({ length: genStats.length }, (_, i) => i + 1),
                datasets: [
                    {
                        label: 'FMin',
                        data: fmin,
                        borderColor: 'red',
                        fill: false,
                    },
                    {
                        label: 'FMax',
                        data: fmax,
                        borderColor: 'green',
                        fill: false,
                    },
                    {
                        label: 'FAvg',
                        data: favg,
                        borderColor: 'blue',
                        fill: false,
                    },
                ],
            };

            const ctx = document.getElementById('myChart').getContext('2d');
            myChart = new Chart(ctx, {
                type: 'line',
                data: chartData,
                options: {
                    responsive: true,
                    plugins: {
                        legend: {
                            position: 'top',
                        },
                    },
                    scales: {
                        x: {
                            title: {
                                display: true,
                                text: 'Generation Index',
                            },
                        },
                        y: {
                            title: {
                                display: true,
                                text: 'Value',
                            },
                        },
                    },
                },
            });

        } else {
            tableBody.innerHTML = "<tr><td colspan='5'>Błąd przy pobieraniu danych</td></tr>";
            console.error("Błąd przy pobieraniu danych");
        }
    } catch (error) {
        tableBody.innerHTML = "<tr><td colspan='5'>Błąd przy pobieraniu danych</td></tr>";
        console.error("Błąd przy pobieraniu danych", error);
    }
}

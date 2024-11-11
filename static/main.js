async function calculate() {
    const a = parseFloat(document.getElementById("a").value);
    const b = parseFloat(document.getElementById("b").value);
    const d = parseFloat(document.getElementById("precision").value); 
    const N = parseInt(document.getElementById("N").value);
    const T = parseInt(document.getElementById("T").value);
    const pk = parseFloat(document.getElementById("pk").value);
    const pm = parseFloat(document.getElementById("pm").value);
    const elite = document.getElementById("elite").checked;

    const requestData = { a: a, b: b, d: d, T: T, N: N, pk: pk, pm: pm, elite: elite };

    // Show loading indicator while fetching data
    const tableBody = document.getElementById("table-body");
    tableBody.innerHTML = "<tr><td colspan='17'>Ładowanie...</td></tr>"; 

    try {
        // Send POST request to the backend
        const response = await fetch("/calculate", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify(requestData)
        });

        if (response.ok) {
            const data = await response.json();
            tableBody.innerHTML = "";

            // Extract generation stats for plotting
            const genStats = data.gen_stats;  // This contains the FMin, FMax, FAvg for each generation
            
            const fmin = [];
            const fmax = [];
            const favg = [];
            
            // Iterate over generation stats and extract the required data
            for (let i = 0; i < genStats.length; i++) {
                fmin.push(genStats[i].f_min);
                fmax.push(genStats[i].f_max);
                favg.push(genStats[i].f_avg);
            }

            // Prepare the chart data
            const chartData = {
                labels: Array.from({ length: genStats.length }, (_, i) => i + 1), // X-axis: generation indexes
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

            // Create the chart
            const ctx = document.getElementById('myChart').getContext('2d');
            const myChart = new Chart(ctx, {
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
            tableBody.innerHTML = "<tr><td colspan='17'>Błąd przy pobieraniu danych</td></tr>";
            console.error("Błąd przy pobieraniu danych");
        }
    } catch (error) {
        tableBody.innerHTML = "<tr><td colspan='7'>Błąd przy pobieraniu danych</td></tr>";
        console.error("Błąd przy pobieraniu danych", error);
    }
}

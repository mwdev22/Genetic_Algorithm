async function calculate() {
    
    // dane z formularza do obliczeń
    const a = parseFloat(document.getElementById("a").value);
    const b = parseFloat(document.getElementById("b").value);
    const d = parseFloat(document.getElementById("precision").value); 
    const N = parseInt(document.getElementById("N").value);

    const requestData = { a: a, b: b, d: d, N: N };

    // ekran ładowania
    const tableBody = document.getElementById("table-body");
    tableBody.innerHTML = "<tr><td colspan='7'>Ładowanie...</td></tr>"; 

    try {
        // pobieranie danych z backendu
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
            console.log(data.g_sum)

            // wypełnianie tabeli danymi

            data.population.forEach((individual) => {
                const row = `<tr id="${individual.id}">
                                <td>${individual.id}</td>
                                <td>${individual.x_real.toFixed(4)}</td>
                                <td>${individual.fx.toFixed(2)}</td>
                                <td>${individual.gx.toFixed(2)}</td>
                             </tr>`;
                tableBody.innerHTML += row;
            });

            selection(data.population, data.g_sum);

            

        } else {
            tableBody.innerHTML = "<tr><td colspan='7'>błąd przy pobieraniu danych</td></tr>";
            console.error("błąd przy pobieraniu danych");
        }
    } catch (error) {
        tableBody.innerHTML = "<tr><td colspan='7'>błąd przy pobieraniu danych</td></tr>";
        console.error("błąd przy pobieraniu danych", error);
    }
}

async function selection(pop, g_sum) {
    const requestData = { pop: pop, g_sum: g_sum };
    console.log(requestData)
    try {
        const response = await fetch("/selection", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify(requestData)
        });
        if (response.ok) {
            const data = await response.json();
            console.log(data)
        } else {
            console.log(response)
        }
    }
    catch (error){
        console.error(error)
    }
}

async function mutation(pop, g_sum) {
    const requestData = {pop: pop, g_sum: g_sum};
    const response = await fetch("/selection", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(requestData)
    });
}

async function crossover(pop, g_sum) {
    const requestData = {pop: pop, g_sum: g_sum};
    const response = await fetch("/selection", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(requestData)
    });
}
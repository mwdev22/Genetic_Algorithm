let prec = 2

async function calculate() {
    
    // dane z formularza do obliczeń
    const a = parseFloat(document.getElementById("a").value);
    const b = parseFloat(document.getElementById("b").value);
    const d = parseFloat(document.getElementById("precision").value); 
    const N = parseInt(document.getElementById("N").value);
    const pk = parseFloat(document.getElementById("pk").value);


    prec = d.toString().length - 2
    const requestData = { a: a, b: b, d: d, N: N };

    // ekran ładowania
    const tableBody = document.getElementById("table-body");
    tableBody.innerHTML = "<tr><td colspan='7'>Ładowanie...</td></tr>"; 

    let expl = document.getElementById('expl')
    expl.innerText = ``


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

            let l = document.getElementById("L");
            l.innerText = data.L

            // wypełnianie tabeli danymi

            data.population.forEach((individual) => {
                const row = `<tr id="${individual.id}">
                                <td>${individual.id}</td>
                                <td>${individual.x_real.toFixed(prec)}</td>
                                <td>${individual.fx.toFixed(prec)}</td>
                                <td>${individual.gx.toFixed(prec)}</td>
                                <td>trwa selekcja...</td>
                                <td>trwa selekcja...</td>
                                <td>trwa selekcja...</td>
                                <td>trwa selekcja...</td>
                                <td>trwa selekcja...</td>
                                <td>trwa krzyżowanie...</td>
                                <td>trwa krzyżowanie...</td>
                                <td>trwa krzyżowanie...</td>
                                <td>trwa krzyżowanie...</td>
                             </tr>`;
                tableBody.innerHTML += row;
            });

            let pop = await selection(data.population, data.g_sum, a, b);
            console.log(pop)
            if (pop != data.population) {
                pop = await crossover(pop, pk);
                console.log(pop)
                // if (pop != data.population) {
                //     // mutation()
                // }
            }


            

        } else {
            tableBody.innerHTML = "<tr><td colspan='7'>błąd przy pobieraniu danych</td></tr>";
            console.error("błąd przy pobieraniu danych");
        }
    } catch (error) {
        tableBody.innerHTML = "<tr><td colspan='7'>błąd przy pobieraniu danych</td></tr>";
        console.error("błąd przy pobieraniu danych", error);
    }
}

async function selection(pop, g_sum, a, b) {
    const requestData = { pop: pop, g_sum: g_sum, a: a, b: b };
    console.log(requestData);
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

            data.population.forEach((individual) => {
                const row = document.getElementById(individual.id);
                if (row) {
                    row.cells[0].textContent = individual.id;
                    row.cells[1].textContent = individual.x_real.toFixed(prec);
                    row.cells[2].textContent = individual.fx.toFixed(prec);
                    row.cells[3].textContent = individual.gx.toFixed(prec);
                    row.cells[4].textContent = individual.p.toFixed(prec);
                    row.cells[5].textContent = individual.q.toFixed(prec);
                    row.cells[6].textContent = individual.r.toFixed(prec);
                    row.cells[7].textContent = individual.x_sel.toFixed(prec);
                    row.cells[8].textContent = individual.x_sel_bin;
                }
            });
            return data.population;
        } else {
            console.log("Błąd przy pobieraniu danych z /selection");
            return pop;
        }
    } catch (error) {
        console.error("Błąd przy przetwarzaniu odpowiedzi /selection", error);
        return pop;
    }
}

async function crossover(pop, pk) {
    const requestData = { pop: pop, pk: pk };
    console.log(requestData)
    try {
        const response = await fetch("/crossover", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify(requestData)
        });

        
        if (response.ok) {
            const data = await response.json();

            console.log(data)

            data.population.forEach((individual) => {
                const row = document.getElementById(individual.id);
                if (row) {
                    row.cells[9].textContent = individual.parent
                    if (individual.pc == 0) {
                        row.cells[10].textContent = '-'
                    } else if (individual.id == data.backup_id) {
                        row.cells[10].textContent = `${individual.pc},${data.backup_pc}`
                        let expl = document.getElementById('expl')
                        expl.innerText = `Osobnik z ID ${data.backup_id} został poddany krzyżowaniu drugi raz, ponieważ liczba rodziców była nieparzysta i ostatni został bez pary.`
                    }
                    else {
                        row.cells[10].textContent = individual.pc
                    }
                    row.cells[11].textContent = individual.child
                    row.cells[12].textContent = individual.new_gen
                }
            });
            return data.population;
        } else {
            console.log("Błąd przy pobieraniu danych z /selection");
            return pop;
        }
    } catch (error) {
        console.error("Błąd przy przetwarzaniu odpowiedzi /selection", error);
        return pop;
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


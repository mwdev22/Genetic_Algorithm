async function calculate() {
    // pobieranie wartości z pól do wprowadzania
    const a = parseFloat(document.getElementById("a").value);
    const b = parseFloat(document.getElementById("b").value);
    const d = parseFloat(document.getElementById("precision").value); 
    const N = parseInt(document.getElementById("N").value);

    const requestData = { a: a, b: b, d: d, N: N };
    
    // request na endpoint kalkulujący wyniki
    const response = await fetch("/calculate", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(requestData)
    });

    if (response.ok) {
        const data = await response.json();
        
        // czyszczenie tabeli i wypełnianie nowymi danymi
        const tableBody = document.getElementById("table-body");
        tableBody.innerHTML = ""; 
        console.log(data)

        data.population.forEach((individual) => {
            const row = `<tr>
                            <td>${individual.id}</td>
                            <td>${individual.x_real.toFixed(4)}</td>
                            <td>${individual.x_int}</td>
                            <td>${individual.bin}</td>
                            <td>${individual.x_new_int}</td>
                            <td>${individual.x_new_real.toFixed(4)}</td>
                            <td>${individual.fx.toFixed(6)}</td>
                         </tr>`;
            tableBody.innerHTML += row;
        });

        // wyświetlanie danych o osobniku z najlepszymi wynikami
        const bestInd = data.best_ind
        console.log(bestInd)
        const resID = document.getElementById("res-id")
        resID.innerText = bestInd.id
        const resVal = document.getElementById("res-val")
        resVal.innerText = bestInd.fx.toFixed(6)
    } else {
        console.error("błąd przy kalkulowaniu danych");
    }
}
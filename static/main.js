async function calculate() {
    
    // dane z formularza do obliczeń
    const a = parseFloat(document.getElementById("a").value);
    const b = parseFloat(document.getElementById("b").value);
    const d = parseFloat(document.getElementById("precision").value); 
    const N = parseInt(document.getElementById("N").value);

    const requestData = { a: a, b: b, d: d, N: N };

    const resID = document.getElementById("res-id");
    const resVal = document.getElementById("res-val");
    const bitCount = document.getElementById("bit-count");

    resID.innerText = '';
    resVal.innerText = '';
    bitCount.innerText = '';

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

            // wypełnianie tabeli danymi

            data.population.forEach((individual) => {
                const row = `<tr id="row-${individual.id}">
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


            // wypełnianie informacji o najlepiej dopasowanym osobniku
            const bestInd = data.best_ind;
            
            resID.innerHTML = `<a href="#row-${bestInd.id}">${bestInd.id}</a>`;
            resVal.innerText = bestInd.fx.toFixed(6);
            bitCount.innerText = data.L;

            // funkcjonalność przewijania do najlepszego osobnika oraz wyróżnianie go
            const previousHighlighted = document.querySelectorAll(".highlighted");
            previousHighlighted.forEach((row) => {
                row.classList.remove("highlighted");
            });

            const bestRow = document.getElementById(`row-${bestInd.id}`);
            if (bestRow) {
                bestRow.classList.add("highlighted");
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


// upiększenie przesuwania do wyniku
document.addEventListener("DOMContentLoaded", () => {
    document.getElementById('res-id').addEventListener('click', function(event) {
        event.preventDefault();  

        const targetID = this.querySelector('a').getAttribute('href').substring(1);  
        const targetElement = document.getElementById(targetID);  

        if (targetElement) {
            targetElement.scrollIntoView({
                behavior: 'smooth',
                block: 'center'
            });
        }
    });
});
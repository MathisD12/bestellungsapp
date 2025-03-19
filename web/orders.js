const table = document.getElementById("orders");
const tableSummaryFood = document.getElementById("summary-food");
const tableSummaryDrinks = document.getElementById("summary-drinks");

async function loadOrders() {
    const response = await fetch("/orders", {
        method: "GET"
    });

    const orders = await response.json();

    console.log("orders:", orders);

    for (const order of orders) {
        const name = document.createElement("td");
        const product = document.createElement("td");
        const veggy = document.createElement("td");
        const drink = document.createElement("td");
        const size = document.createElement("td");

        name.innerText = order.Name;
        product.innerText = order.Product.Name;
        veggy.innerText = order.Product.Veggy ? 'veggy' : 'nicht veggy';
        
        if (order.Drink) {
            drink.innerText = order.Drink.Name;
            size.innerText = order.Drink.Size >= 2 ? 'groß' : 'klein';
        }

        const row = document.createElement("tr");
        row.appendChild(name);
        row.appendChild(product);
        row.appendChild(veggy);
        row.appendChild(drink);
        row.appendChild(size);

        table.appendChild(row);
    }
}

async function loadSummary() {
    const response = await fetch("/orders/summary", {
        method: "GET"
    });

    const summary = await response.json();

    console.log("summary:", summary);

    for (const product in summary.Products) {
        const key = document.createElement("td");
        const count = document.createElement("td");

        key.innerText = product;
        count.innerText = summary.Products[product].toString();

        const row = document.createElement("tr");
        row.appendChild(key);
        row.appendChild(count);

        tableSummaryFood.appendChild(row);
    }

    for (const product in summary.Drinks) {
        const key = document.createElement("td");
        const count = document.createElement("td");

        key.innerText = product;
        count.innerText = summary.Drinks[product].toString();

        const row = document.createElement("tr");
        row.appendChild(key);
        row.appendChild(count);

        tableSummaryDrinks.appendChild(row);
    }
}

loadOrders();
loadSummary();


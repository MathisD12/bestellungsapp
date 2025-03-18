const button = document.getElementById('button-absenden');
const inputname = document.getElementById('input-name');
const inputEssen = document.getElementById('input-essen');
const inputtrinken = document.getElementById('input-trinken');
const checkboxVeggy = document.getElementById('checkbox-veggy');
const selectSize = document.getElementById('select-size');
const fehler = document.getElementById('fehler');

button.onclick = async () => {
    const name = inputname.value;
    const essen = inputEssen.value;
    const trinken = inputtrinken.value;
    const veggy = checkboxVeggy.checked;
    const size = Number(selectSize.value);

    const order = {
        name,
        product: {
            name: essen,
            veggy
        }
    };

    if (trinken) {
        order.drink = {
            name: trinken,
            size
        };
    }

    console.log(`Bestellung: ${name} ${essen} ${trinken}`);

    const response = await fetch('/orders', {
        method: 'POST',
        body: JSON.stringify(order)
    });
    
    if (!response.ok) {
        fehler.innerText = await response.text();
    }
}

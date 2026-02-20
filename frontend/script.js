function showText() {
    const input = document.getElementById("login");
    let value = input.value;      // то, что ввёл пользователь
    if (value !== "") {
        const textElement = document.querySelector('.messages');
        const newElement = document.createElement('div');
        if (value.includes("ХУЙ")) {
            newElement.className = 'hui';
            newElement.innerHTML = value;
            textElement.before(newElement);
            console.log(newElement);
        } else {
            newElement.className = 'mes';
            newElement.innerHTML = value;
            textElement.before(newElement);
            console.log(newElement);
        }
    }
    const box = document.getElementById("chatic");
    box.scrollTop = box.scrollHeight;
}
// function clear() {
//     // получаем все элементы с классом 'mes'
//     const clearElements = document.querySelectorAll('.mes');

//     // перебираем и удаляем каждый
//     clearElements.forEach(el => el.remove());
// }
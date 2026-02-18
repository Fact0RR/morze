function showText() {
    const input = document.getElementById("login");
    let value = input.value;      // то, что ввёл пользователь
    const textElement = document.querySelector('.messages');
    const newElement = document.createElement('div');
    newElement.className = 'mes';
    newElement.innerHTML = value;
    textElement.before(newElement);
    console.log(newElement);
}
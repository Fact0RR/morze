showHistory()
// Создаём асинхронную функцию для запроса данных
async function getData() {

  try {
    // Отправляем GET-запрос и ЖДЁМ, пока сервер ответит
    const response = await fetch(
      'http://78.24.219.132:8081/api/messages?limit=20&offset=0&contact=1'
    );

    // Читаем тело ответа и превращаем JSON в JS-объект
    const messages = await response.json();

    // Используем полученные данные (пока просто выводим)
    return messages;

  } catch (error) {
    // Если сервер не ответил или данные сломались — ловим ошибку
    console.error('Ошибка:', error);
  }
}

async function showHistory() {
  const messages = await getData();
  for (let i = messages.length-1; i>=0; i--) {
    const textElement = document.querySelector('.messages');
    const newElement = document.createElement('div');
    console.log(messages[i].data, messages[i].created_at, messages[i].user_id);
    if (messages[i].user_id === 1) {
      console.log("привет");
      newElement.className = 'mes';
      newElement.innerHTML = messages[i].data;
      textElement.before(newElement);
    } else {
      console.log("хуй");
      newElement.className = 'hui';
      newElement.innerHTML = messages[i].data;
      textElement.before(newElement);
    }
  }
}


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
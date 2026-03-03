// Создаём асинхронную функцию для запроса данных
showHistory()
async function getData() {
  const GetWay = '/api/messages';
  const Limit = '?limit=';
  const Offset = '&offset=';
  const Contact = '&contact=';
  let IDlimit = 50;
  let IDoffset = 0;
  let IDcontact = 1;
  try {
    // Отправляем GET-запрос и ЖДЁМ, пока сервер ответит
    const response = await fetch(envProtocol + envIPadres + envPort + GetWay + Limit + IDlimit + Offset + IDoffset + Contact + IDcontact);

    // Читаем тело ответа и превращаем JSON в JS-объект
    const messages = await response.json();

    // Используем полученные данные (пока просто выводим)
    return messages;

  } catch (error) {
    // Если сервер не ответил или данные сломались — ловим ошибку
    console.error('Ошибка:', error);
  }
}

async function postData(value, userID) {
  const PostWay = '/api/message';

  const message = {
    contact_id: 1,
    user_id: userID,
    data: value
  };

  try {
    // Отправляем GET-запрос и ЖДЁМ, пока сервер ответит
    const response = await fetch(envProtocol + envIPadres + envPort + PostWay, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(message)
    });

    // Читаем тело ответа и превращаем JSON в JS-объект
    const message_id = await response.json();

    // Используем полученные данные (пока просто выводим)
    return message_id;

  } catch (error) {
    // Если сервер не ответил или данные сломались — ловим ошибку
    console.error('Ошибка:', error);
  }
}

// нужно сделать эту функцию асинхронной походу const messages = await getData();
async function showText() {
  const input = document.getElementById("login");
  let value = input.value;      // то, что ввёл пользователь
  let inputRadio = document.querySelector('input[name="idOneOrTwo"]:checked');
  let userID = 2;
  if (inputRadio.value === "One") {
    userID = 1;
  }
  if (value !== "") {
    const messages = await postData(value, userID);
    const textElement = document.querySelector('.messages');
    const newElement = document.createElement('div');
    if (userID === 1) {
      newElement.className = 'mesFirstUser';
      newElement.innerHTML = value;
      textElement.before(newElement);
    } else {
      newElement.className = 'mesSecondUser';
      newElement.innerHTML = value;
      textElement.before(newElement);
    }
  }
  const box = document.getElementById("chatic");
  box.scrollTop = box.scrollHeight;
}

async function showHistory() {
  const messages = await getData();
  for (let i = messages.length - 1; i >= 0; i--) {
    messages[i].created_at = new Date(messages[i].created_at);
  }
  messages.sort((a, b) => a.created_at - b.created_at);
  for (let i = 0; i < messages.length; i++) {
    const textElement = document.querySelector('.messages');
    const newElement = document.createElement('div');
    if (messages[i].user_id === 1) {
      newElement.className = 'mesFirstUser';
      newElement.innerHTML = messages[i].data;
      textElement.before(newElement);
    } else {
      newElement.className = 'mesSecondUser';
      newElement.innerHTML = messages[i].data;
      textElement.before(newElement);
    }
  }
}

window.onload = () => {
  const chat = document.getElementById('chatic');
  setTimeout(() => {
    chat.scrollTop = chat.scrollHeight;
  }, 50);
};
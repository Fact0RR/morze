// Создаём асинхронную функцию для запроса данных
const Protocol = 'http://';
const IPadres = '78.24.219.132';
const Port = ':8081';
const GetWay = '/api/messages'
const PostWay = '/api/message'
const Limit = '?limit='
const Offset = '&offset='
const Contact = '&contact='
let IDlimit = '50';
let IDoffset = '0';
let IDcontact = '1';
let GetUrl = Protocol + IPadres + Port + GetWay + Limit + IDlimit + Offset + IDoffset + Contact + IDcontact;
let PostUrl = Protocol + IPadres + Port + PostWay;
showHistory()
async function getData() {

  try {
    // Отправляем GET-запрос и ЖДЁМ, пока сервер ответит
    const response = await fetch(GetUrl);

    // Читаем тело ответа и превращаем JSON в JS-объект
    const messages = await response.json();

    // Используем полученные данные (пока просто выводим)
    return messages;

  } catch (error) {
    // Если сервер не ответил или данные сломались — ловим ошибку
    console.error('Ошибка:', error);
  }
}

async function postData(value) {

  const message = {
    contact_id: 1,
    user_id: 1,
    data: value
  };

  try {
    // Отправляем GET-запрос и ЖДЁМ, пока сервер ответит
    const response = await fetch(PostUrl, {
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

async function showHistory() {
  const messages = await getData();
  for (let i = messages.length - 1; i >= 0; i--) {
    const textElement = document.querySelector('.messages');
    const newElement = document.createElement('div');
    if (messages[i].user_id === 1) {
      newElement.className = 'mesFirstUser';
      newElement.innerHTML = messages[i].data;
      textElement.before(newElement);
      console.log(messages[i].id);
    } else {
      newElement.className = 'mesSecondUser';
      newElement.innerHTML = messages[i].data;
      textElement.before(newElement);
      console.log(messages[i].id);
    }
  }
}

// нужно сделать эту функцию асинхронной походу const messages = await getData();
async function showText() {
  const input = document.getElementById("login");
  let value = input.value;      // то, что ввёл пользователь
  if (value !== "") {
    const messages = await postData(value);
    const textElement = document.querySelector('.messages');
    const newElement = document.createElement('div');
    newElement.className = 'mesFirstUser';
    newElement.innerHTML = value;
    textElement.before(newElement);
    console.log(messages.message_id)
  }
  const box = document.getElementById("chatic");
  box.scrollTop = box.scrollHeight;
}
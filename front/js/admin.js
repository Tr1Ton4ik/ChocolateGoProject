const tabs = document.querySelectorAll('nav a');
const sections = document.querySelectorAll('main > section');

tabs.forEach((tab, index) => {
  tab.addEventListener('click', (e) => {
    e.preventDefault();
    sections.forEach(section => section.style.display = 'none');
    sections[index].style.display = 'block';
  });
});

function getCookie(name) {
  let cookie = document.cookie.split('; ').find(row => row.startsWith(name + '='));
  return cookie ? cookie.split('=')[1] : null;
}

if (getCookie('isAuthenticated') !== "true") {
  window.location.href = '/login';
}
document.getElementById('admin-form').addEventListener('submit', (event) => {
  event.preventDefault();

  const password = document.getElementById('password').value;
  const confirmPassword = document.getElementById('confirm-password').value;

  if (password !== confirmPassword) {
    alert('Пароли не совпадают. Пожалуйста, проверьте еще раз.');
    return;
  }

  // Если пароли совпадают, отправляем форму
  event.target.submit();
});

// Обработчик клика для кнопки выхода
document.getElementById('logout-btn').addEventListener('click', (event) => {
  event.preventDefault();

  // Отправляем запрос на сервер, чтобы выйти из аккаунта
  fetch('/logout', {
    method: 'POST',
    credentials: 'include'
  })
    .then(response => {
      if (response.ok) {
        // Если выход успешен, перенаправляем пользователя на страницу входа
        document.cookie = 'isAuthenticated=false; path=/';
        window.location.href = '/login';
      } else {
        alert('Произошла ошибка при выходе из аккаунта. Пожалуйста, попробуйте еще раз.');
      }
    })
    .catch(error => {
      console.error('Ошибка при выходе из аккаунта:', error);
      alert('Произошла ошибка при выходе из аккаунта. Пожалуйста, попробуйте еще раз.');
    });
});
document.getElementById('admin-form').addEventListener('submit', (event) => {
  event.preventDefault();

  const password = document.getElementById('password').value;
  const confirmPassword = document.getElementById('confirm-password').value;

  if (password !== confirmPassword) {
    alert('Пароли не совпадают. Пожалуйста, проверьте еще раз.');
    return;
  }

  // Если пароли совпадают, отправляем форму
  event.target.submit();
});
async function searchProduct() {
  const query = document.getElementById('search-input').value;
  const response = await fetch(`/search_product?name=${encodeURIComponent(query)}`);
  const results = await response.json();
  displayResults(results);
}

// Функция для отображения результатов поиска
function displayResults(results) {
  const resultsContainer = document.getElementById('search-results');
  resultsContainer.innerHTML = ''; // очищаем предыдущие результаты
  if (results.length > 0) {
      results.forEach(product => {
          const div = document.createElement('div');
          div.textContent = product.name + `(${product.category})`; // Предполагаем, что у товара есть свойство 'name'
          resultsContainer.appendChild(div);
      });
  } else {
      resultsContainer.textContent = 'Нет результатов';
  }
}
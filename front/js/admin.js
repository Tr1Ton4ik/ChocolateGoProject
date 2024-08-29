const tabs = document.querySelectorAll('nav a');
const sections = document.querySelectorAll('main > section');

tabs.forEach((tab, index) => {
    tab.addEventListener('click', (e) => {
        e.preventDefault();
        sections.forEach(section => section.style.display = 'none');
        sections[index].style.display = 'block';
    });
});
document.querySelectorAll('form button[type="submit"]').forEach(button => {
    button.addEventListener('click', (event) => {
        event.preventDefault();

        // Показываем подтверждение действия
        const confirmed = confirm('Вы действительно хотите выполнить это действие?');

        if (confirmed) {
            // Если пользователь подтвердил, отправляем форму
            event.target.closest('form').submit();
        }
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
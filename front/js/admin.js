const tabs = document.querySelectorAll('nav a');
const sections = document.querySelectorAll('main > section');

tabs.forEach((tab, index) => {
    tab.addEventListener('click', (e) => {
        e.preventDefault();
        sections.forEach(section => section.style.display = 'none');
        sections[index].style.display = 'block';
    });
});
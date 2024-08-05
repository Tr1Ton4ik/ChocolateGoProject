let cart = JSON.parse(localStorage.getItem('cart')) || [];

const addToCartButtons = document.querySelectorAll('.add-to-cart');
addToCartButtons.forEach(button => {
  button.addEventListener('click', addToCart);
});

function addToCart(event) {
  const button = event.currentTarget;
  const productInfo = {
    name: button.dataset.name,
    price: parseFloat(button.dataset.price)
  };
  cart.push(productInfo);
  localStorage.setItem('cart', JSON.stringify(cart));
  updateCartLink();
}

function updateCartLink() {
  const cartLink = document.querySelector('nav a[href="cart.html"]');
  cartLink.textContent = "Cart (" + cart.length + ")";
}
function redirectToProductPage(element) {
  const productName = element.getAttribute('alt');
  window.location.href = `product.html?name=${encodeURIComponent(productName)}`;
}
updateCartLink();
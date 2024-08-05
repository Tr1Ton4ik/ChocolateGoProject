let cart = JSON.parse(localStorage.getItem('cart')) || [];

function displayCart() {
  const cartItemsContainer = document.getElementById('cart-items');
  cartItemsContainer.innerHTML = '';

  const cartMap = new Map();

  cart.forEach(item => {
    if (cartMap.has(item.name)) {
      cartMap.get(item.name).quantity++;
      cartMap.get(item.name).total += item.price;
    } else {
      cartMap.set(item.name, { quantity: 1, total: item.price, price: item.price, image: item.image });
    }
  });

  let totalPrice = 0;

  cartMap.forEach((item, name) => {
    const cartItem = document.createElement('div');
    cartItem.classList.add('cart-item');

    const productImage = document.createElement('img');
    productImage.src ='/front/img/'+name+'.jpeg'; // Assuming you have the image URL stored in item.image
    productImage.alt = name;
    productImage.classList.add('product-image');
    productImage.onclick = () => redirectToProductPage(name); // Add click event handler

    const productName = document.createElement('p');
    productName.textContent = name;

    const productQuantity = document.createElement('p');
    productQuantity.textContent = `x${item.quantity}`;

    const productPrice = document.createElement('p');
    productPrice.textContent = `$${item.total.toFixed(2)}`;

    cartItem.appendChild(productImage);
    cartItem.appendChild(productName);
    cartItem.appendChild(productQuantity);
    cartItem.appendChild(productPrice);

    cartItemsContainer.appendChild(cartItem);

    totalPrice += item.total;
  });

  document.getElementById('cart-total').textContent = totalPrice.toFixed(2);
}

function redirectToProductPage(productName) {
  window.location.href = `product.html?name=${encodeURIComponent(productName)}`;
}

const clearCartButton = document.getElementById('clear-cart');
clearCartButton.addEventListener('click', () => {
  cart = [];
  localStorage.setItem('cart', JSON.stringify(cart));
  displayCart();
});

displayCart();

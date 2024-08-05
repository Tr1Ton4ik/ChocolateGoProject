function displayProductDetails() {
    const urlParams = new URLSearchParams(window.location.search);
    const productId = urlParams.get('id');
  
    fetch(`/api/products/${productId}`)
      .then(response => response.json())
      .then(product => {
        const productImage = document.querySelector('.product-image');
        const productName = document.querySelector('.product-name');
        const productDescription = document.querySelector('.product-description');
        const productPrice = document.querySelector('.product-price');
        const addToCartButton = document.querySelector('.add-to-cart');
  
        productImage.src = product.image;
        productName.textContent = product.name;
        productDescription.textContent = product.description;
        productPrice.textContent = `$${product.price.toFixed(2)}`;
        addToCartButton.dataset.name = product.name;
        addToCartButton.dataset.price = product.price;
      })
      .catch(error => {
        console.error('Error fetching product details:', error);
      });
  }
  
  window.addEventListener('load', displayProductDetails);
function generateProduct(name) {
  const product = document.createElement("div");
  product.classList.add("product");
  product.id = `product-${name}`;
  const title = document.createElement("h2");
  title.innerText = name;
  const deleteButton = document.createElement("button");
  deleteButton.innerText = "Delete";
  deleteButton.id = name;
  deleteButton.onclick = deleteItem;
  product.appendChild(title);
  product.appendChild(deleteButton);
  return product;
}
async function fetchData() {
  try {
    const res = await fetch("http://localhost:8080/get");
    const data = await res.json();
    console.log(data);
    if (data === null) {
      return;
    }
    const productsDOM = document.getElementById("products");
    for (i = 0; i < data.length; i++) {
      productsDOM.appendChild(generateProduct(data[i].Name));
    }
  } catch (error) {
    console.error("Error fetching data:", error);
  }
}
fetchData();
async function deleteItem(e) {
  name = e.target.id;
  const res = await fetch("/delete", {
    method: "POST",
    mode: "cors",
    headers: {
      "Content-Type": "application/x-www-form-urlencoded",
    },
    body: `product_name=${name}`,
  });
  document
    .getElementById(`products`)
    .removeChild(document.getElementById(`product-${name}`));
}
async function newItem(e) {
  e.preventDefault();
  value = document.getElementById("item").value;
  const res = await fetch("/add", {
    method: "POST",
    mode: "cors",
    headers: {
      "Content-Type": "application/x-www-form-urlencoded",
    },
    body: `product_name=${value}`,
  });
  if (res.status == 200) {
    document.getElementById("products").appendChild(generateProduct(value));
  }
}

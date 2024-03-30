import { FormEvent, useEffect, useRef, useState } from "react";
import "./App.css";

interface Product {
  ID: string;
  Name: string;
}

function App() {
  const [data, setData] = useState<undefined | null | Array<Product>>(
    undefined,
  );
  const [error, setError] = useState<string | null>(null);
  const NewProductInputRef = useRef<null | HTMLInputElement>(null);
  console.log(`This application is running in ${import.meta.env.MODE} mode.`);
  let domain: string;
  if (import.meta.env.MODE === "development") {
    domain = `http://${window.location.host.split(":")[0]}:8030`;
  } else {
    domain = "";
  }
  console.log(domain);
  useEffect(() => {
    async function FetchGet() {
      try {
        const res = await fetch(`${domain}/api/get`);
        if (!res.ok) {
          console.log(res);
          throw new Error("Failed to fetch data");
        }
        const jsonData = await res.json();
        setData(jsonData);
      } catch (error) {
        if (error instanceof Error) {
          setError(error.message);
        }
      }
    }

    FetchGet();
  }, [domain]);
  async function deleteAllProducts() {
    if (!data) {
      return "This function cannot be called. There must be data present and loaded.";
    }
    for (let i = 0; i < data.length; i++) {
      deleteProduct(data[i].ID, "all");
    }
    setData(null);
  }
  async function deleteProduct(id: string, origin: string) {
    if (!origin) {
      origin = "one";
    }
    const res = await fetch(`${domain}/api/delete`, {
      method: "POST",
      mode: "cors",
      headers: {
        "Content-Type": "application/x-www-form-urlencoded",
      },
      body: `ID=${id}`,
    });
    if (res.status === 200 && origin != "all") {
      const newData = data?.filter((v) => v.ID != id);
      setData(newData);
    }
  }
  async function NewProductHandler(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const value = NewProductInputRef.current?.value;
    const res = await fetch(`${domain}/api/add`, {
      method: "POST",
      mode: "cors",
      headers: {
        "Content-Type": "application/x-www-form-urlencoded",
      },
      body: `product_name=${value}`,
    });
    const jsonData: Product = await res.json();
    if (res.status === 200) {
      let newData;
      if (!data) {
        newData = [jsonData];
      } else {
        newData = [...(data as Product[]), jsonData];
      }
      setData(newData);
    }
  }
  function NewProductField() {
    return (
      <form
        className="newProduct"
        onSubmit={(e) => {
          NewProductHandler(e);
        }}
      >
        <input
          className="newProductInput"
          placeholder="Item"
          ref={NewProductInputRef}
          name="item-name"
        ></input>
        <button type="submit" className="newProductButton">
          Add
        </button>
      </form>
    );
  }
  if (error) {
    return <div>Error: {error}</div>;
  }
  if (data) {
    return (
      <>
        <button
          onClick={(e) => {
            e.stopPropagation();
            e.preventDefault();
            deleteAllProducts();
          }}
        >
          Delete All
        </button>
        {data.map((v) => (
          <div
            className="Item"
            key={v.ID}
            onClick={(e) => {
              e.stopPropagation();
              e.preventDefault();
              deleteProduct(v.ID, "one");
            }}
          >
            <span key={v.ID}>- {v.Name}</span>
          </div>
        ))}
        <NewProductField />
      </>
    );
  }
  if (data === undefined) {
    return <div>Loading...</div>;
  }
  return <NewProductField />;
}

export default App;

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
    domain = "http://localhost:8030";
  } else {
    domain = "";
  }
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
  async function deleteProduct(id: string) {
    const res = await fetch(`${domain}/api/delete`, {
      method: "POST",
      mode: "cors",
      headers: {
        "Content-Type": "application/x-www-form-urlencoded",
      },
      body: `ID=${id}`,
    });
    if (res.status === 200) {
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
      const newData = [...(data as Product[]), jsonData];
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
          placeholder="product"
          ref={NewProductInputRef}
          name="item-name"
        ></input>
        <button type="submit" className="newProductButton">
          submit
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
        {data.map((v) => (
          <div
            className="Item"
            key={v.ID}
            onClick={(e) => {
              e.stopPropagation();
              e.preventDefault();
              deleteProduct(v.ID);
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

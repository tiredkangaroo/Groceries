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
  useEffect(() => {
    async function FetchGet() {
      try {
        const res = await fetch(`/api/get`);
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
  }, []);
  async function deleteProduct(id: string) {
    const res = await fetch(`/api/delete`, {
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
    const res = await fetch(`/api/add`, {
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
        onSubmit={(e) => {
          NewProductHandler(e);
        }}
      >
        <input placeholder="product" ref={NewProductInputRef}></input>
        <button type="submit">submit</button>
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
          <div>
            <input
              className="cb"
              type="checkbox"
              onInput={(e) => {
                e.stopPropagation();
                e.preventDefault();
                deleteProduct(v.ID);
              }}
            ></input>
            <span key={v.ID}>{v.Name}</span>
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

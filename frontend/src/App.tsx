import { useEffect, useRef, useState } from "react";
import "./App.css";
interface Item {
  ID: string;
  DateCreated: string;
  Name: string;
}

function resolveHost(): string {
  let host = "";
  if (import.meta.env.MODE === "development") {
    host = window.location.protocol + "//" + window.location.hostname + ":8030";
  } else {
    host = window.location.protocol + "//" + window.location.host;
  }
  return host;
}

async function retrieveJSON(uRL: string): Promise<any> {
  const resp = await fetch(resolveHost() + uRL);
  const data = await resp.json();
  return data;
}

async function retrieveItems(): Promise<[Array<Item>, string]> {
  const data = await retrieveJSON("/api/items");

  if (data.error != null) {
    return [[], data.error];
  }
  var items: Array<Item> = [];
  for (let i = 0; i < data.items.length; i++) {
    const item = data.items[i];
    const newItem: Item = {
      ID: item.id,
      DateCreated: item.date_created,
      Name: item.name,
    };
    items.push(newItem);
  }
  return [items, ""];
}

async function deleteItem(
  items: Array<Item>,
  setItems: React.Dispatch<React.SetStateAction<Array<Item>>>,
  itemID: string,
) {
  const newItems = items.filter((item) => item.ID != itemID);
  setItems(newItems);
  fetch(resolveHost() + "/api/item/" + itemID, {
    method: "DELETE",
  });
}

async function newItem(
  items: Array<Item>,
  setItems: React.Dispatch<React.SetStateAction<Array<Item>>>,
  newItemRef: React.RefObject<HTMLInputElement>,
) {
  const value = newItemRef.current?.value;
  fetch(resolveHost() + "/api/item", {
    method: "POST",
    body: value,
  }).then((resp) => {
    resp.json().then((data) => {
      if (data.error != null) {
        console.error(data.error);
        return;
      }
      setItems([
        ...items,
        {
          ID: data.item.id,
          DateCreated: data.item.date_created,
          Name: data.item.name,
        },
      ]);
    });
  });
}

function App() {
  const [error, setError] = useState("");
  const [items, setItems] = useState<Array<Item>>([]);
  const newItemRef = useRef<HTMLInputElement>(null);
  useEffect(() => {
    async function retrieve() {
      const [itemsResponse, err] = await retrieveItems();
      console.log(itemsResponse);
      if (err != "") {
        setError(err);
        return;
      }
      setItems(itemsResponse);
    }
    retrieve();
  }, []);
  return (
    <div className="big">
      <h1>{error}</h1>
      <div className="items">
        <h1 className="items-title"> checklist application </h1>
        {items.map((item) => {
          return (
            <div
              key={item.ID}
              className="item"
              onClick={() => {
                deleteItem(items, setItems, item.ID);
              }}
            >
              <p className="item-name">{item.Name}</p>
            </div>
          );
        })}
        <input
          className="items-input"
          placeholder="new item"
          ref={newItemRef}
        ></input>
        <button
          className="items-submit"
          type="button"
          onClick={() => newItem(items, setItems, newItemRef)}
        >
          Submit
        </button>
      </div>
    </div>
  );
}

export default App;

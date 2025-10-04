import React from "react";
import styles from "./SearchDone.module.scss";
import Search from "@/modules/Search/Search";
import { Link } from "react-router-dom";

export default function PageSearch() {
  return (
    <div>
      <label>Страница поиска</label>
      <Search/>
      <button>Поиск</button>
    </div>
  )
}

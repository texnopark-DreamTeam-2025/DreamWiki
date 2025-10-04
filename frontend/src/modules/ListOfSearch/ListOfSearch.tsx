import React, { useState } from 'react';
import search from "Public/static/icons/iconSearch.svg";
import styles from "./ListOfSearch.module.scss";


export default function ListOfSearch() {
  let word = useState("")
  return (
    <div className={styles.searchForm}>
      <img className={styles.searchSearchForm__icon} src={search}></img>
      <input
        type="text"
        placeholder="Поиск"
        className={styles.searchForm__input}
        onChange={ (event: { target: { value: string } }) => {debouncedSetSearch(event.target.value)}}
      ></input>
    </div>
  );
}

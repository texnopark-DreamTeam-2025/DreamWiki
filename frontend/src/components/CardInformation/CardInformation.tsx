import styles from "./CardInformation.module.scss";

export default function CardInformation({
  title,
  description,
}: {
  title: string;
  description: string;
}) {
  return (
    <div className={styles.cardInformation}>
      <h2>{title}</h2>
      <p>{description}</p>
    </div>
  );
}

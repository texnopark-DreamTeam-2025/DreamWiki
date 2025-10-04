export default function CardInformation(title: string, description: string) {
  return (
    <div className="card-information">
      <h2>{title}</h2>
      <p>{description}</p>
    </div>
  );
}

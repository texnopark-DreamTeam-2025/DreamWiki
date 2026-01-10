import { Button, Card, Text, Flex } from "@gravity-ui/uikit";
import { ArrowUpRightFromSquare } from "@gravity-ui/icons";
import ReactMarkdown from "react-markdown";

export default function CardInformation({
  title,
  description,
  onDetailsClick,
  yandexWikiLink,
}: {
  title: string;
  description: string;
  onDetailsClick?: () => void;
  yandexWikiLink?: string;
}) {
  return (
    <Card>
      <Flex direction="column" gap="4" className="p-2">
        <Text variant="subheader-3">{title}</Text>
        <Flex direction="column">
          <ReactMarkdown>{description}</ReactMarkdown>
        </Flex>
        <Flex direction="row">
          <Button onClick={onDetailsClick} view="flat">
            Подробнее
          </Button>
          {yandexWikiLink && (
            <Button
              onClick={(e) => {
                e.stopPropagation();
                window.open(yandexWikiLink, "_blank");
              }}
              view="flat"
            >
              <Flex alignItems="center">
                <ArrowUpRightFromSquare />
                Яндекс Вики
              </Flex>
            </Button>
          )}
        </Flex>
      </Flex>
    </Card>
  );
}

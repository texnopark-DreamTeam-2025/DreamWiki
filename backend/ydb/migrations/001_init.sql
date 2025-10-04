CREATE TABLE Page (
    page_id UUID,
    content TEXT NOT NULL,
    PRIMARY KEY (page_id)
);

CREATE TABLE Paragraph (
    paragraph_id UUID NOT NULL,
    page_id UUID NOT NULL,
    content TEXT NOT NULL,
    PRIMARY KEY (paragraph_id)
);

import * as preact from "preact";
import { ComponentChildren } from "preact";

import styles from "./Paragraph.module.css";

type Props = {
  children: ComponentChildren;
};

const Paragraph = ({ children }: Props) => {
  return <p className={styles.paragraph}>{children}</p>;
};

export default Paragraph;

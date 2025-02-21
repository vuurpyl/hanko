import * as preact from "preact";
import { ComponentChildren } from "preact";

import styles from "./Content.module.css";

type Props = {
  children: ComponentChildren;
};

const Content = ({ children }: Props) => {
  return <section className={styles.content}>{children}</section>;
};

export default Content;

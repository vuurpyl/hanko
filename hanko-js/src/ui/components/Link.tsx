import * as preact from "preact";
import { ComponentChildren, FunctionalComponent } from "preact";
import cx from "classnames";

import styles from "./Link.module.css";

export type Props = {
  children?: ComponentChildren;
  onClick?: (event: Event) => void;
  disabled?: boolean;
  hidden?: boolean;
};

const Link: FunctionalComponent<Props> = ({
  children,
  onClick,
  disabled,
  hidden,
}: Props) => {
  return (
    <a
      onClick={onClick}
      hidden={hidden}
      className={cx(styles.link, disabled ? styles.disabled : null)}
    >
      {children}
    </a>
  );
};

export default Link;

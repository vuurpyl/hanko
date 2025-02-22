import * as preact from "preact";
import { ComponentChildren, toChildArray } from "preact";

import styles from "./Form.module.css";

type Props = {
  onSubmit: (event: Event) => void;
  children: ComponentChildren;
};

const Form = ({ onSubmit, children }: Props) => {
  return (
    <form onSubmit={onSubmit}>
      <ul className={styles.ul}>
        {toChildArray(children).map((child, index) => (
          <li className={styles.li} key={index}>
            {child}
          </li>
        ))}
      </ul>
    </form>
  );
};

export default Form;

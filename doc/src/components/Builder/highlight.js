import React from "react"
import Prism, {defaultProps} from "prism-react-renderer"
import clsx from "clsx";
import styles from "./styles.module.css"


const Highlight = ({code, language = "yaml"}) => {
    return (
        <div>
            <Prism
                {...defaultProps}
                language={language}
                code={code}
            >
                {({className, style, tokens, getLineProps, getTokenProps}) => (
                    <pre className={clsx(styles.highlight, className)} style={style}>
          {tokens.map((line, i) => (
              <div key={i} {...getLineProps({line, key: i})}>
                  {line.map((token, key) => (
                      <span key={key} {...getTokenProps({token, key})} />
                  ))}
              </div>
          ))}
        </pre>
                )}
            </Prism>
        </div>
    )
}
export default Highlight
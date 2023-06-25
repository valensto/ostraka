import React, {useCallback, useEffect, useState} from "react"
import useWindowWidth from "../../hooks/useWindowWidth"
import clsx from "clsx"
import Highlight from "./Highlight"
import Chevron from "./chevron"
import styles from "./styles.module.css"
import {Section} from "../Section";

const S = [60, -300, -700, -1050]
const M = [100, -300, -680, -1020]
const L = [140, -300, -700, -1100]

const getTopByIndex = (m, index) => {
    const scale = {
        1: (m[0] ?? 0),
        2: (m[1] ?? 0),
        3: (m[2] ?? 0),
        4: (m[3] ?? 0),
    }

    return scale[index] ?? 0
}

const eventType = `event_type:
  format: json
  fields:
    - name: customerId
      data_type: string
      required: true
    - name: orderNumber
      data_type: int
      required: true
    - name: orderStatus
      data_type: string
      required: true
    - name: nonRequiredField
      data_type: string`

const Middlewares = `middlewares:
  cors:
    default:
      allowed_origins:
        - http://localhost:3000
      allowed_methods:
        - POST
      allowed_headers:
        - Authorization
      allow_credentials: true
      max_age: 3600

  auth:
    default:
      type: token
      params:
        token: 2dc7929e5b589cb7861bcae19e13ad96
        query_param: token`

const inputFields = `inputs:
  - name: webhook-orders
    source: webhook
    params:
      endpoint: /webhook/orders
      auth: default
    decoder:
      format: json
      mappers:
        - source: o_customer_id
          target: customerId
        - ...`

const mergeQuery = `outputs:
  - name: sse-orders-completed
    destination: sse
    params:
      endpoint: /sse/orders/completed
      auth: default
      cors: default
    encoder:
      format: json
    condition:
      operator: or
      conditions:
        - field: orderStatus
          operator: eq
          value: "completed"
        - field: orderStatus
          operator: eq
          value: "pending"`

export const Builder = () => {
    const [top, setTop] = useState(S)
    const [index, setIndex] = useState(1)
    const windowWidth = useWindowWidth()

    const handleClick = useCallback((index) => {
        setIndex(index)
    }, [])
    const handleUpClick = useCallback(() => {
        setIndex(Math.max(index - 1, 1))
    }, [index])
    const handleDownClick = useCallback(() => {
        setIndex(Math.min(index + 1, 4))
    }, [index])

    useEffect(() => {
        if (windowWidth != null && windowWidth < 622) {
            setTop(S)
            return
        }

        if (windowWidth != null && windowWidth < 800) {
            setTop(M)
            return
        }

        setTop(L)
    }, [windowWidth])

    return (
        <Section fullWidth odd center noGap>
            <section
                className={clsx(
                    styles.section,
                    styles["section--inner"],
                    styles["section--center"],
                    styles["section--showcase"],
                )}
            >
                <h2
                    className={clsx(
                        styles.section__title,
                        styles["section__title--wide"],
                        "text--center",
                    )}
                >
                    Your perfect workflow in minutes
                </h2>

                <p
                    className={clsx(
                        styles.section__subtitle,
                        styles["section__subtitle--narrow"],
                        "text--center",
                    )}
                >
                    No code required. Just write YAML to create the workflow which match your existing infrastructure
                    requirements
                </p>

                <div className={styles.showcase}>
                    <div className={styles.showcase__inner}>
                        <div
                            className={clsx(styles.showcase__chevron)}
                            onClick={handleUpClick}
                            style={{visibility: index === 1 ? "hidden" : "visible"}}
                        >
                            <Chevron/>
                        </div>
                        <div className={clsx(styles.showcase__left)}>
                            <div
                                className={clsx(
                                    styles.showcase__offset,
                                    styles[`showcase__${index}`],
                                )}
                                style={{top: getTopByIndex(top, index)}}
                            >
                                <Highlight code={eventType}/>
                                <Highlight code={Middlewares}/>
                                <Highlight code={inputFields}/>
                                <Highlight code={mergeQuery}/>
                            </div>
                        </div>
                        <div
                            className={clsx(
                                styles.showcase__chevron,
                                styles["showcase__chevron--bottom"],
                            )}
                            onClick={handleDownClick}
                            style={{visibility: index === 4 ? "hidden" : "visible"}}
                        >
                            <Chevron/>
                        </div>
                        <div className={styles.showcase__right}>
                            <div
                                className={clsx(styles.showcase__button, {
                                    [styles["showcase__button--active"]]: index === 1,
                                })}
                                onClick={() => handleClick(1)}
                            >
                                <h3 className={styles.showcase__header}>
                                    Define your Event type
                                </h3>
                                <p className={styles.showcase__description}>
                                    Define your event type and its fields
                                </p>
                            </div>

                            <div
                                className={clsx(styles.showcase__button, {
                                    [styles["showcase__button--active"]]: index === 2,
                                })}
                                onClick={() => handleClick(2)}
                            >
                                <h3 className={styles.showcase__header}>
                                    Create your middlewares
                                </h3>
                                <p className={styles.showcase__description}>
                                    Create auth and cors middlewares to secure your endpoints
                                </p>
                            </div>

                            <div
                                className={clsx(styles.showcase__button, {
                                    [styles["showcase__button--active"]]: index === 3,
                                })}
                                onClick={() => handleClick(3)}
                            >
                                <h3 className={styles.showcase__header}>
                                    Configure your inputs
                                </h3>
                                <p className={styles.showcase__description}>
                                    Add as many inputs as you want to your workflow
                                </p>
                            </div>
                            <div
                                className={clsx(styles.showcase__button, {
                                    [styles["showcase__button--active"]]: index === 4,
                                })}
                                onClick={() => handleClick(4)}
                            >
                                <h3 className={styles.showcase__header}>
                                    Configure your outputs
                                </h3>
                                <p className={styles.showcase__description}>
                                    Add as many outputs as you want to your workflow and add a condition to filter the
                                    data
                                </p>
                            </div>
                        </div>
                    </div>
                </div>
            </section>
        </Section>
    )
}
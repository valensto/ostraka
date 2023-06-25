import React from "react"
import clsx from "clsx"
import {Section} from "../Section";
import styles from "./styles.module.css"

export const WhyCards = () => (
    <Section fullWidth>
        <Section noGap>
            <Section.Title level={3} size="small" center>
                Why another dispatcher?
            </Section.Title>

            <div
                className={clsx(
                    styles.section__footer,
                    styles["section__footer--feature-cards"],
                )}
            >
                {[
                    {
                        header: "Developer-centric approach",
                        content:
                            "Ostraka is designed with developers in mind, providing a powerful event dispatching solution that offers flexibility and control over event processing.",
                    },

                    {
                        header: "Seamless integration",
                        content:
                            "Ostraka seamlessly integrates into existing systems and architectures, allowing developers to leverage their current technology stack and infrastructure.",
                    },

                    {
                        header: "Customization and extensibility",
                        content:
                            "Ostraka empowers developers to tailor event dispatching to their specific needs, with support for a wide range of event sources and destinations.",
                    },

                    {
                        header: "Harnessing software technologies",
                        content:
                            "Ostraka leverages proven software technologies and best practices, ensuring reliability, scalability, and performance in event processing.",
                    },

                    {
                        header: "Simplified event management",
                        content:
                            "With Ostraka, developers can easily manage and orchestrate events across different systems and services, enabling efficient communication and synchronization.",
                    },

                    {
                        header: "Empowering the developer community",
                        content:
                            "Ostraka contributes to the open-source developer community by providing a valuable and flexible event dispatching tool, fostering collaboration and innovation.",
                    },
                ].map(({ header, content }, index) => (
                    <div key={index} className={styles.feature}>
                        <h3 className={styles.feature__header}>{header}</h3>
                        <p className={styles.feature__content}>{content}</p>
                    </div>
                ))}
            </div>
        </Section>
    </Section>
)
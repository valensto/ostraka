import React from 'react';
import {Section} from "../Section";
import styles from "./styles.module.css";
import Highlight from "../Builder/highlight";
import {workflow} from "./workflow"

function Editor(props) {
    return (
        <div className={styles.preview}>
            <Section.Subtitle className={styles.previewHeader}>
                Try QuestDB demo in your browser
            </Section.Subtitle>

            <div className={styles.editor}>
                <div className={styles.code}>
                    <Highlight code={workflow} />
                </div>
            </div>
        </div>
    );
}

export default Editor;
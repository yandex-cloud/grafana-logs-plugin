import React, { ChangeEvent, useState } from 'react';
import { Button, Field, IconButton, InlineField, InlineFieldRow, InlineSwitch, Input } from '@grafana/ui';

interface DerivedFiedProps {
    name: string
    value: string
    onChange: (value: string) => void;
    onDelete: () => void;
}
export const DerivedField: React.FC<DerivedFiedProps> = (props) => {
    const { name, onChange, onDelete, value } = props;

    const [state, setState] = useState({ curValue: value, });
    const { curValue } = state

    return (
        <InlineField label={name} transparent>
            <Input
                key={name}
                suffix={<IconButton
                    name='times'
                    size="sm"
                    onClick={onDelete}
                />}
                placeholder='mustache template...'
                value={curValue}
                onChange={(event: ChangeEvent<HTMLInputElement>) => {
                    setState({ ...state, curValue: event.target.value });
                }}
                onKeyDown={(e) => {
                    if (e.key === 'Enter') {
                        e.preventDefault();
                    }
                }}
                onKeyUp={(e) => {
                    e.preventDefault();
                    if (e.key === 'Enter') {
                        onChange(curValue);
                    }
                }}
                onBlur={(e) => {
                    onChange(curValue)
                }}
            />
        </InlineField>
    );
}

interface AddDerivedFieldsProps {
    onConfirm: (name: string) => void;
}
export const AddDerivedField: React.FC<AddDerivedFieldsProps> = (props) => {
    const { onConfirm } = props;

    const [state, setState] = useState({ name: "", });

    const { name } = state;
    const confirm = () => {
        if (name !== '') {
            onConfirm(name);
        }
        setState({ ...state, name: "" });
    }

    return (
        <InlineField>
            <Input
                placeholder='add field...'
                suffix={<IconButton
                    name='plus'
                    variant='primary'
                    onClick={confirm}
                />}
                value={name}
                onChange={(event: ChangeEvent<HTMLInputElement>) => {
                    setState({ ...state, name: event.target.value });
                }}
                onKeyDown={(e) => {
                    // onKeyDown is triggered before onKeyUp, triggering submit behaviour on Enter press if this component
                    // is used inside forms. Moving onKeyboardAdd callback here doesn't work since text input is not captured in onKeyDown
                    if (e.key === 'Enter') {
                        e.preventDefault();
                    }
                }}
                onKeyUp={(e) => {
                    e.preventDefault();
                    if (e.key === 'Enter') {
                        confirm();
                    }
                }}
            />
        </InlineField>
    );
}

export interface DerivedFieldData {
    name: string;
    value: string;
}

interface DerivedFieldsProps {
    fields: DerivedFieldData[];
    onChange: (fields: DerivedFieldData[]) => void;
}
export const DerivedFields: React.FC<DerivedFieldsProps> = (props) => {
    const { fields, onChange } = props;

    return (
        <InlineFieldRow>
            {fields.map((v, i) => (
                <DerivedField
                    key={v.name}
                    name={v.name}
                    value={v.value}
                    onChange={(v) => {
                        const newFields = [...fields]
                        newFields[i].value = v
                        onChange(newFields)
                    }}
                    onDelete={() => {
                        const newFields = [...fields]
                        newFields.splice(i, 1)
                        onChange(newFields)
                    }}
                />
            ))}
            <AddDerivedField
                onConfirm={(name) => {
                    onChange([...fields, { name, value: "" }])
                }}
            />
        </InlineFieldRow>
    );
}



export interface DerivedLinkData {
    name?: string;
    title?: string;
    url?: string;
    targetBlank?: boolean;
}

interface DerivedLinkProps {
    value: DerivedLinkData;
    onChange: (value: DerivedLinkData) => void;
    onDelete: () => void;
}

export const DerivedLink: React.FC<DerivedLinkProps> = (props) => {
    const { value, onChange, onDelete } = props;

    return (
        <Field >
            <InlineFieldRow>
                <InlineField label="Field">
                    <Input placeholder="field name..."
                        value={value.name}
                        onChange={(event: ChangeEvent<HTMLInputElement>) => {
                            onChange({ ...value, name: event.target.value })
                        }}
                    />
                </InlineField>
                <InlineField label="Title">
                    <Input placeholder="link title..."
                        value={value.title}
                        onChange={(event: ChangeEvent<HTMLInputElement>) => {
                            onChange({ ...value, title: event.target.value })
                        }}
                    />
                </InlineField>
                <InlineField label="URL" grow>
                    <Input type="url" placeholder="http://example.com/${__value.raw}"
                        value={value.url || ""}
                        onChange={(event: ChangeEvent<HTMLInputElement>) => {
                            onChange({ ...value, url: event.target.value })
                        }}
                    />
                </InlineField>
                <InlineSwitch
                    label="In new tab"
                    showLabel={true}
                    value={value.targetBlank || false}
                    onChange={(event: ChangeEvent<HTMLInputElement>) => {
                        onChange({ ...value, targetBlank: event.target.checked })
                    }}
                />
                <Button icon='times' variant="destructive" onClick={onDelete} />
            </InlineFieldRow>
        </Field>
    );
}


interface DerivedLinksProps {
    links: DerivedLinkData[];
    onChange: (links: DerivedLinkData[]) => void;
}
export const DerivedLinks: React.FC<DerivedLinksProps> = (props) => {
    const { links, onChange } = props;

    return (
        <div>
            {links.map((v, i) => (
                <DerivedLink key={i}
                    value={v}
                    onChange={(newValue) => {
                        const newLinks = [...links]
                        newLinks.splice(i, 1, newValue)
                        onChange(newLinks)
                    }}
                    onDelete={() => {
                        const newLinks = [...links]
                        newLinks.splice(i, 1)
                        onChange(newLinks)
                    }}
                />
            ))}
            <Field>
                <Button
                    icon='plus'
                    variant="secondary"
                    onClick={() => onChange([...links, {}])}
                >
                    Add URL
                </Button>
            </Field>
        </div>
    );
}

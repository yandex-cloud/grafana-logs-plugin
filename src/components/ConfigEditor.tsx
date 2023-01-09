import React, { ChangeEvent, PureComponent } from 'react';
import { Field, FieldSet, Input, SecretTextArea } from '@grafana/ui';
import { DataSourcePluginOptionsEditorProps } from '@grafana/data';
import { LoggingSourceOptions, LoggingSecureJsonData, defaultSourceOptions, LoggingDerivedLink } from '../types';
import { defaults } from 'lodash';
import { DerivedLinks, DerivedLinkData } from './Derived';

interface Props extends DataSourcePluginOptionsEditorProps<LoggingSourceOptions> { }

interface State { }

export class ConfigEditor extends PureComponent<Props, State> {
  onEndpointChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onOptionsChange, options } = this.props;
    const jsonData = {
      ...options.jsonData,
      apiEndpoint: event.target.value,
    };
    onOptionsChange({ ...options, jsonData });
  };

  onFolderChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onOptionsChange, options } = this.props;
    const jsonData = {
      ...options.jsonData,
      folderId: event.target.value,
    };
    onOptionsChange({ ...options, jsonData });
  };

  onLinksChange = (links: DerivedLinkData[]) => {
    const { onOptionsChange, options } = this.props;
    const jsonData = {
      ...options.jsonData,
      derivedLinks: links.map((v): LoggingDerivedLink => ({
        field: v.name || "",
        title: v.title || "",
        url: v.url || "",
        targetBlank: v.targetBlank || false
      })),
    };
    onOptionsChange({ ...options, jsonData });
  };

  // Secure field (only sent to the backend)
  onAPIKeyChange = (event: ChangeEvent<HTMLTextAreaElement>) => {
    const { onOptionsChange, options } = this.props;
    onOptionsChange({
      ...options,
      secureJsonData: {
        apiKeyJson: event.target.value,
      },
    });
  };

  onResetAPIKey = () => {
    const { onOptionsChange, options } = this.props;
    onOptionsChange({
      ...options,
      secureJsonFields: {
        ...options.secureJsonFields,
        apiKeyJson: false,
      },
      secureJsonData: {
        ...options.secureJsonData,
        apiKeyJson: '',
      },
    });
  };

  render() {
    const { options } = this.props;
    const { secureJsonFields } = options;
    const jsonData = defaults(options.jsonData, defaultSourceOptions) as LoggingSourceOptions;
    const secureJsonData = (options.secureJsonData || {}) as LoggingSecureJsonData;
    const links: DerivedLinkData[] = (jsonData.derivedLinks || []).map((item) => ({ name: item.field, title: item.title, url: item.url, targetBlank: item.targetBlank }))

    return (
      <div>
        <FieldSet label="SDK config">
          <Field label="API endpoint">
            <Input
              onChange={this.onEndpointChange}
              value={jsonData.apiEndpoint}
              placeholder="yandex cloud api endpoint <host>:<port>"
            />
          </Field>
          <Field label="Folder ID">
            <Input
              onChange={this.onFolderChange}
              value={jsonData.folderId || ""}
              placeholder="folder for log groups search"
            />
          </Field>
        </FieldSet>

        <FieldSet label="Derived DataLinks">
          <Field>
            <div>Use derived data links to add custom links to derived log fields.</div>
          </Field>
          <DerivedLinks links={links} onChange={this.onLinksChange} />
        </FieldSet>

        <FieldSet label="Secret config">
          <Field label="API Key">
            <SecretTextArea
              rows={8}
              isConfigured={(secureJsonFields && secureJsonFields.apiKeyJson) as boolean}
              value={secureJsonData.apiKeyJson || ''}
              placeholder="place full json key file content here"
              onReset={this.onResetAPIKey}
              onChange={this.onAPIKeyChange}
            />
          </Field>
        </FieldSet>
      </div >
    );
  }
}

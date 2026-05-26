//#region src/integration.d.ts
interface DevToolsIntegrationOptions {
  config: unknown;
}
declare function DevToolsIntegration(options: DevToolsIntegrationOptions): {
  name: string;
};
declare function runDevTools(builder: unknown): Promise<void>;
//#endregion
export { DevToolsIntegration, DevToolsIntegrationOptions, runDevTools };
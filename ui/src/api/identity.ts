export type Identity =
  | HostIdentity
  | KubernetesIdentity

export type HostIdentity = {
  kind: "host"
  hostname: string
  process?: string
}

export type KubernetesIdentity = {
  kind: "kubernetes"
  apiVersion: string
  component: string
  name: string
  namespace: string
  cluster?: string
}

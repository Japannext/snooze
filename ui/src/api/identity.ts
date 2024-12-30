export type Identity =
  | HostIdentity
  | KubernetesIdentity

export type HostIdentity = {
  kind: "host"
  host: string
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

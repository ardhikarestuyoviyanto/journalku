export default function checkPermission(permission) {
  const auth = localStorage.getItem('auth')
  if (!auth) return false

  const data = JSON.parse(auth)
  const allPermission = data.user?.currentCompany?.permission || []

  return allPermission.includes(permission)
}

'use client'

import { useState, useCallback } from 'react'
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card"
import { AlertCircle, Copy, Link } from 'lucide-react'

export default function LinkShortener() {
  const [original, setOriginal] = useState('')
  const [short, setShort] = useState('')
  const [error, setError] = useState('')
  const [isLoading, setIsLoading] = useState(false)
  const [response, setResponse] = useState<null | { short: string; original: string }>(null)

  const handleSubmit = useCallback(async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    setResponse(null)
    setIsLoading(true)

    if (!original) {
      setError('Please enter a URL')
      setIsLoading(false)
      return
    }

    try {
      new URL(original)
    } catch {
      setError('Please enter a valid URL')
      setIsLoading(false)
      return
    }

    const formData = new FormData()
    formData.append('original', original)
    if (short) {
      formData.append('short', short)
    }

    try {
      const res = await fetch('http://localhost:1323/shorten', {
        method: 'POST',
        body: formData,
      })

      if (!res.ok) {
        throw new Error('Failed to shorten URL')
      }

      const data = await res.json()
      setResponse(data)
    } catch (err) {
      setError('An error occurred while shortening the URL')
      console.error(err)
    } finally {
      setIsLoading(false)
    }
  }, [original, short])

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text)
  }

  return (
    <Card className="w-full max-w-md mx-auto">
      <CardHeader>
        <CardTitle>Link Shortener</CardTitle>
        <CardDescription>Enter a long URL to get a shortened version</CardDescription>
      </CardHeader>
      <CardContent>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="original">Long URL</Label>
            <Input
              id="original"
              type="url"
              placeholder="https://example.com/very/long/url"
              value={original}
              onChange={(e) => setOriginal(e.target.value)}
              required
            />
          </div>
          <div className="space-y-2">
            <Label htmlFor="short">Custom short code (optional)</Label>
            <Input
              id="short"
              type="text"
              placeholder="custom-code"
              value={short}
              onChange={(e) => setShort(e.target.value)}
            />
          </div>
          {error && (
            <div className="text-red-500 flex items-center gap-2">
              <AlertCircle className="h-4 w-4" />
              <span>{error}</span>
            </div>
          )}
          <Button type="submit" className="w-full" disabled={isLoading}>
            {isLoading ? 'Shortening...' : 'Shorten URL'}
          </Button>
        </form>
      </CardContent>
      {response && (
        <CardFooter className="flex flex-col items-start gap-4">
          <div className="w-full space-y-2">
            <Label>Shortened URL</Label>
            <div className="flex w-full items-center gap-2">
              <Input value={`${window.location.origin}/${response.short}`} readOnly className="flex-grow" />
              <Button variant="outline" size="icon" onClick={() => copyToClipboard(`${window.location.origin}/${response.short}`)}>
                <Copy className="h-4 w-4" />
                <span className="sr-only">Copy shortened URL</span>
              </Button>
            </div>
          </div>
          <div className="w-full space-y-2">
            <Label>Server Response</Label>
            <pre className="bg-muted p-2 rounded-md overflow-x-auto w-full">
              <code>{JSON.stringify(response, null, 2)}</code>
            </pre>
          </div>
        </CardFooter>
      )}
    </Card>
  )
}